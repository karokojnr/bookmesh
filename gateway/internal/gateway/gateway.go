package gateway

import (
	"context"

	pb "github.com/karokojnr/bookmesh-shared/proto"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrder(ctx context.Context, orderId, customerId string) (*pb.Order, error)
}
