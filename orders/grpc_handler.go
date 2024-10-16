package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/karokojnr/bookmesh-shared/api"
	"github.com/karokojnr/bookmesh-shared/broker"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	svc         OrdersService
	amqpChannel *amqp.Channel
}

func NewGrpcHandler(grpcServer *grpc.Server, svc OrdersService, amqpChannel *amqp.Channel) {
	handler := &grpcHandler{
		svc:         svc,
		amqpChannel: amqpChannel,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	/// show gRPC error handling
	/// return nil, fmt.Errorf("not implemented")

	log.Printf("Creating order for customer %v", req)

	books, err := h.svc.ValidateOrder(ctx, req)
	if err != nil {
		return nil, err
	}

	o, err := h.svc.CreateOrder(ctx, req, books)
	if err != nil {
		return nil, err
	}

	mOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	q, err := h.amqpChannel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	h.amqpChannel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         mOrder,
		DeliveryMode: amqp.Persistent,
	})

	return o, nil
}

func (h *grpcHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	log.Printf("Getting order for customer %v", req)
	return h.svc.GetOrder(ctx, req)

}

func (h *grpcHandler) UpdateOrder(ctx context.Context, p *pb.Order) (*pb.Order, error) {
	return h.svc.UpdateOrder(ctx, p)
}
