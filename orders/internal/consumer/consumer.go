package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/karokojnr/bookmesh-orders/internal/types"
	"github.com/karokojnr/bookmesh-shared/broker"
	pb "github.com/karokojnr/bookmesh-shared/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type consumer struct {
	service types.OrdersService
}

func NewConsumer(svc types.OrdersService) *consumer {
	return &consumer{
		service: svc,
	}
}
func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(q.Name, "", broker.OrderPaidEvent, false, nil)
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
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}

			_, err := c.service.UpdateOrder(context.Background(), o)
			if err != nil {
				log.Printf("failed to update order: %v", err)

				if err := broker.HandleRetry(ch, &d); err != nil {
					log.Printf("Error handling retry: %v", err)
				}

				continue
			}

			messageSpan.AddEvent("order.updated")
			messageSpan.End()

			log.Println("Order has been updated from AMQP")
			d.Ack(false)
		}
	}()

	<-forever
}
