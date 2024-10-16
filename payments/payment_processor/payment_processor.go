package paymentprocessor

import (
	pb "github.com/karokojnr/bookmesh-shared/api"
)

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
