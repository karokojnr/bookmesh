package main

import (
	"context"

	"github.com/karokojnr/bookmesh-payments/gateway"
	paymentprocessor "github.com/karokojnr/bookmesh-payments/payment_processor"
	pb "github.com/karokojnr/bookmesh-shared/api"
)

type service struct {
	processor paymentprocessor.PaymentProcessor
	gateway   gateway.OrdersGateway
}

func NewService(processor paymentprocessor.PaymentProcessor, gateway gateway.OrdersGateway) *service {
	return &service{
		processor: processor,
		gateway:   gateway,
	}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	/// connect to payment process
	link, err := s.processor.CreatePaymentLink(o)
	if err != nil {
		return "", err
	}

	/// Update order with the payment link
	err = s.gateway.UpdateOrderWithPaymentLink(ctx, o.OrderId, link)
	if err != nil {
		return "", err
	}

	return link, nil
}
