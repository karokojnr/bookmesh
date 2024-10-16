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

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever = make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Received message %s: ", d.Body)

			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				log.Println("Failed to unmarshal message: ", err)
				continue
			}

			paymentLink, err := c.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Println("Failed to create payment: ", err)
				continue
			}

			log.Println("Payment link created: ", paymentLink)
		}
	}()
	<-forever
}
