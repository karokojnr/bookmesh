package service

import (
	"context"
	"testing"

	"github.com/karokojnr/bookmesh-payments/gateway"
	inmemoryprocessor "github.com/karokojnr/bookmesh-payments/internal/payment_processor/in_memory_processor"
	inMemoryRegistry "github.com/karokojnr/bookmesh-shared/discovery/in_memory"
	pb "github.com/karokojnr/bookmesh-shared/proto"
)

func TestService(t *testing.T) {
	p := inmemoryprocessor.NewInMemoryProcessor()
	registry := inMemoryRegistry.NewRegistry()
	gateway := gateway.NewGateway(registry)
	svc := NewService(p, gateway)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &pb.Order{})

		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}

		if link == "" {
			t.Errorf("Expected a payment link but got an empty string")
		}
	})
}
