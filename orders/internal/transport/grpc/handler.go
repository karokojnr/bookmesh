package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/karokojnr/bookmesh-shared/broker"
	pb "github.com/karokojnr/bookmesh-shared/proto"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Book) (*pb.Order, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	UpdateOrder(context.Context, *pb.Order) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Book, error)
}

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

	q, err := h.amqpChannel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	tr := otel.Tracer("amqp")
	amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - produce - %s", q.Name))
	defer messageSpan.End()

	books, err := h.svc.ValidateOrder(amqpContext, req)
	if err != nil {
		return nil, err
	}

	o, err := h.svc.CreateOrder(amqpContext, req, books)
	if err != nil {
		return nil, err
	}

	mOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	/// Inject headers for tracing (context propagation)
	headers := broker.InjectAMQPHeaders(amqpContext)

	h.amqpChannel.PublishWithContext(amqpContext, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         mOrder,
		DeliveryMode: amqp.Persistent,
		Headers:      headers,
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
