package main

import (
	"context"

	paymentprocessor "github.com/karokojnr/bookmesh-payments/payment_processor"
	pb "github.com/karokojnr/bookmesh-shared/api"
)

type service struct {
	processor paymentprocessor.PaymentProcessor
}

func NewService(processor paymentprocessor.PaymentProcessor) *service {
	return &service{
		processor: processor,
	}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	/// connect to payment process
	link, err := s.processor.CreatePaymentLink(o)
	if err != nil {
		return "", err
	}

	/// Update order with the payment link

	return link, nil
}
