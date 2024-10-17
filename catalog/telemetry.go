package main

import (
	"context"
	"fmt"

	pb "github.com/karokojnr/bookmesh-shared/api"
	"go.opentelemetry.io/otel/trace"
)

type TelemetryMiddleware struct {
	next CatalogService
}

func NewTelemetryMiddleware(next CatalogService) CatalogService {
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
