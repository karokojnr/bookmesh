package transport

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/karokojnr/bookmesh-orders/internal/types"
	"github.com/karokojnr/bookmesh-shared/broker"
	pb "github.com/karokojnr/bookmesh-shared/proto"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"

	amqp "github.com/rabbitmq/amqp091-go"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	svc         types.OrdersService
	amqpChannel *amqp.Channel
}

func NewGrpcHandler(grpcServer *grpc.Server, svc types.OrdersService, amqpChannel *amqp.Channel) {
	handler := &grpcHandler{
		svc:         svc,
		amqpChannel: amqpChannel,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	// show gRPC error handling
	// return nil, fmt.Errorf("not implemented")

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

	// Inject headers for tracing (context propagation)
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
	return h.svc.GetOrder(ctx, req)

}

func (h *grpcHandler) UpdateOrder(ctx context.Context, p *pb.Order) (*pb.Order, error) {
	return h.svc.UpdateOrder(ctx, p)
}
