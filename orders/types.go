package main

import (
	"context"

	pb "github.com/karokojnr/bookmesh-shared/api"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Book, error)
}

type OrdersStore interface {
	Create(context.Context) error
}
