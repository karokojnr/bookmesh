package paymentprocessor

import (
	pb "github.com/karokojnr/bookmesh-shared/proto"
)

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
