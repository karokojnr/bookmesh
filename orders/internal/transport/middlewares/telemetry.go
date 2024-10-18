package middlewares

import (
	"context"
	"fmt"

	"github.com/karokojnr/bookmesh-orders/internal/types"
	pb "github.com/karokojnr/bookmesh-shared/proto"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryMiddleware struct {
	next types.OrdersService
}

func NewTelemetryMiddleware(next types.OrdersService) *TelemetryMiddleware {
	return &TelemetryMiddleware{next}
}

func (s *TelemetryMiddleware) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, books []*pb.Book) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CreateOrder: %v", p))
	return s.next.CreateOrder(ctx, p, books)

}

func (s *TelemetryMiddleware) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("GetOrder: %v", req))
	return s.next.GetOrder(ctx, req)
}

func (s *TelemetryMiddleware) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("UpdateOrder: %v", o))
	return s.next.UpdateOrder(ctx, o)
}

func (s *TelemetryMiddleware) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Book, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("ValidateOrder: %v", req))
	return s.next.ValidateOrder(ctx, req)
}
