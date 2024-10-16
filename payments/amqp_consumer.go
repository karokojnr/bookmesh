package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/karokojnr/bookmesh-shared/api"
	"github.com/karokojnr/bookmesh-shared/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	service PaymentService
}

func NewConsumer(svc PaymentService) *consumer {
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
			log.Printf("Received message %s: ", d.Body)

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

			log.Println("Payment link created: ", paymentLink)
			d.Ack(false)
		}
	}()
	<-forever
}
