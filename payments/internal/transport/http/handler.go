package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	shared "github.com/karokojnr/bookmesh-shared"
	"github.com/karokojnr/bookmesh-shared/broker"
	pb "github.com/karokojnr/bookmesh-shared/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
	"go.opentelemetry.io/otel"
)

var endpointStripeSecret = shared.EnvString("ENDPOINT_STRIPE_SECRET", "")

type PaymentHandler struct {
	amqpChannel *amqp.Channel
}

func NewHttpPaymentHandler(amqpChannel *amqp.Channel) *PaymentHandler {
	return &PaymentHandler{amqpChannel}
}

func (h *PaymentHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("POST/webhook", h.handleCheckoutWebhook)
}

func (h *PaymentHandler) handleCheckoutWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	fmt.Fprintf(os.Stdout, "Received webhook: %s\n", body)

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), endpointStripeSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	if event.Type == "checkout.session.completed" {

		var cs stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &cs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if cs.PaymentStatus == "paid" {
			log.Printf("Payment for Checkout Session %s succeeded!", cs.ID)
			// publish message to RabbitMQ
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			orderId := cs.Metadata["orderId"]
			customerId := cs.Metadata["customerId"]

			o := &pb.Order{
				OrderId:     orderId,
				CustomerId:  customerId,
				Status:      "paid",
				PaymentLink: "",
			}

			marshalledOrder, err := json.Marshal(o)
			if err != nil {
				log.Fatal(err.Error())
			}

			/// Trace
			tr := otel.Tracer("amqp")
			amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - produce - %s", broker.OrderPaidEvent))
			defer messageSpan.End()

			/// Inject headers for tracing (context propagation)
			headers := broker.InjectAMQPHeaders(amqpContext)

			h.amqpChannel.PublishWithContext(amqpContext, broker.OrderPaidEvent, "", false, false, amqp.Publishing{
				ContentType:  "application/json",
				Body:         marshalledOrder,
				DeliveryMode: amqp.Persistent,
				Headers:      headers,
			})

			log.Println("Message published order.paid")
		}

	}

	w.WriteHeader(http.StatusOK)
}
