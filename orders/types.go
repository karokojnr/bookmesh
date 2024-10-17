package main

import (
	"context"

	pb "github.com/karokojnr/bookmesh-shared/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Book) (*pb.Order, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	UpdateOrder(context.Context, *pb.Order) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Book, error)
}

type OrdersStore interface {
	Create(context.Context, Order) (primitive.ObjectID, error)
	Get(ctx context.Context, id, customerID string) (*Order, error)
	Update(ctx context.Context, id string, o *pb.Order) error
}

type Order struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	CustomerId  string             `bson:"customer_id,omitempty"`
	Status      string             `bson:"status,omitempty"`
	PaymentLink string             `bson:"payment_link,omitempty"`
	Books       []*pb.Book         `bson:"books,omitempty"`
}

func (o *Order) ToProto() *pb.Order {
	return &pb.Order{
		OrderId:     o.Id.Hex(),
		CustomerId:  o.CustomerId,
		Status:      o.Status,
		PaymentLink: o.PaymentLink,
	}
}
