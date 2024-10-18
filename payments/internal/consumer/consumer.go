package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/karokojnr/bookmesh-payments/internal/types"
	"github.com/karokojnr/bookmesh-shared/broker"
	pb "github.com/karokojnr/bookmesh-shared/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type consumer struct {
	service types.PaymentService
}

func NewConsumer(svc types.PaymentService) *consumer {
	return &consumer{
		service: svc,
	}
}

func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever = make(chan struct{})

	go func() {
		for d := range msgs {

			/// Extract the headers from the message
			ctx := broker.ExtractAMQPHeader(context.Background(), d.Headers)

			tr := otel.Tracer("amqp")
			_, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - consumer - %s", q.Name))

			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				d.Nack(false, false)
				log.Println("Failed to unmarshal message: ", err)
				continue
			}

			paymentLink, err := c.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Printf("failed to create payment: %v", err)

				if err := broker.HandleRetry(ch, &d); err != nil {
					log.Printf("Error handling retry: %v", err)
				}

				d.Nack(false, false)

				continue
			}

			messageSpan.AddEvent(fmt.Sprintf("payment.created: %s", paymentLink))
			messageSpan.End()

			log.Printf("Payment link created %s", paymentLink)
			d.Ack(false)
		}
	}()
	<-forever
}
