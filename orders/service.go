package main

import (
	"context"
	"log"

	shared "github.com/karokojnr/bookmesh-shared"
	pb "github.com/karokojnr/bookmesh-shared/api"
)

type service struct {
	store OrdersStore
}

func NewService(store OrdersStore) *service {
	return &service{store: store}
}

func (s *service) CreateOrder(ctx context.Context) error {
	return nil
}

func (s *service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) error {
	if len(req.Books) == 0 {
		return shared.ErrNoBooks
	}

	mergedBooks := mergeBooksQuantities(req.Books)
	log.Println("Merged books:", mergedBooks)

	/// validate with the catalog service
	return nil
}

func mergeBooksQuantities(books []*pb.BookWithQuantity) []*pb.BookWithQuantity {
	merged := make(map[string]*pb.BookWithQuantity)
	for _, b := range books {
		if _, ok := merged[b.BookId]; !ok {
			merged[b.BookId] = b
			continue
		}
		merged[b.BookId].Quantity += b.Quantity
	}
	var res []*pb.BookWithQuantity
	for _, b := range merged {
		res = append(res, b)
	}
	return res
}
