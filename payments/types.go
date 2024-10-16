package main

import (
	"context"

	pb "github.com/karokojnr/bookmesh-shared/api"
)

type PaymentService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
