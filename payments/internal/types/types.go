package types

import (
	"context"

	pb "github.com/karokojnr/bookmesh-shared/proto"
)

type PaymentService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
