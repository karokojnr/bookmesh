package middlewares

import (
	"context"
	"time"

	"github.com/karokojnr/bookmesh-orders/internal/types"
	pb "github.com/karokojnr/bookmesh-shared/proto"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	next types.OrdersService
}

func NewLoggingMiddleware(next types.OrdersService) types.OrdersService {
	return &LoggingMiddleware{next}
}

func (s *LoggingMiddleware) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Book) (*pb.Order, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("CreateOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.CreateOrder(ctx, p, items)
}

func (s *LoggingMiddleware) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("GetOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.GetOrder(ctx, p)
}

func (s *LoggingMiddleware) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("UpdateOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.UpdateOrder(ctx, o)
}

func (s *LoggingMiddleware) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Book, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("ValidateOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.ValidateOrder(ctx, p)
}
