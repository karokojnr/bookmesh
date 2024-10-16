package inmemoryprocessor

import (
	pb "github.com/karokojnr/bookmesh-shared/api"
)

type InMemoryProcessor struct{}

func NewInMemoryProcessor() *InMemoryProcessor {
	return &InMemoryProcessor{}
}

func (i *InMemoryProcessor) CreatePaymentLink(o *pb.Order) (string, error) {
	return "test-payment-link", nil
}
