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

	o, err := h.svc.CreateOrder(ctx, req)

	// o := &pb.Order{
	// 	OrderId: "1",
	// 	Books:   req.Books,
	// }

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
