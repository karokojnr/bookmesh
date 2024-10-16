package main

import (
	"context"
	"errors"

	pb "github.com/karokojnr/bookmesh-shared/api"
)

var orders = make([]*pb.Order, 0)

type storage struct {
	// mongo db
}

func NewStorage() *storage {
	return &storage{}
}

func (s *storage) Create(ctx context.Context, req *pb.CreateOrderRequest, books []*pb.Book) (string, error) {
	id := "order_1"
	orders = append(orders, &pb.Order{
		OrderId:     id,
		CustomerId:  req.CustomerId,
		Status:      "pending",
		Books:       books,
		PaymentLink: "",
	})
	return id, nil
}

func (s *storage) Get(ctx context.Context, orderId, customerId string) (*pb.Order, error) {
	for _, o := range orders {
		if o.OrderId == orderId && o.CustomerId == customerId {
			return o, nil
		}
	}
	return nil, errors.New("order not found")
}

func (s *storage) Update(ctx context.Context, orderId string, order *pb.Order) error {
	for i, o := range orders {
		if o.OrderId == orderId {
			orders[i].Status = order.Status
			orders[i].PaymentLink = order.PaymentLink
			return nil
		}
	}
	return errors.New("order not found")
}
