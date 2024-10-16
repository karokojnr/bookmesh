package main

import (
	"context"
	"testing"

	inmemoryprocessor "github.com/karokojnr/bookmesh-payments/payment_processor/in_memory_processor"
	pb "github.com/karokojnr/bookmesh-shared/api"
)

func TestService(t *testing.T) {
	p := inmemoryprocessor.NewInMemoryProcessor()
	svc := NewService(p)

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
