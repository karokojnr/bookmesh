package middleware

import (
	"context"
	"fmt"

	"github.com/karokojnr/bookmesh-catalog/internal/types"
	pb "github.com/karokojnr/bookmesh-shared/proto"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryMiddleware struct {
	next types.CatalogService
}

func NewTelemetryMiddleware(next types.CatalogService) types.CatalogService {
	return &TelemetryMiddleware{next}
}

func (s *TelemetryMiddleware) GetBooks(ctx context.Context, ids []string) ([]*pb.Book, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("GetItems: %v", ids))

	return s.next.GetBooks(ctx, ids)
}

func (s *TelemetryMiddleware) CheckIfBookIsInCatalog(ctx context.Context, p []*pb.BookWithQuantity) (bool, []*pb.Book, error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(fmt.Sprintf("CheckIfItemAreInStock: %v", p))

	return s.next.CheckIfBookIsInCatalog(ctx, p)
}
