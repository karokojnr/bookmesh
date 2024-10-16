package main

import (
	"context"

	pb "github.com/karokojnr/bookmesh-shared/api"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Book) (*pb.Order, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Book, error)
}

type OrdersStore interface {
	Create(context.Context, *pb.CreateOrderRequest, []*pb.Book) (string, error)
	Get(ctx context.Context, orderId, customerId string) (*pb.Order, error)
}
