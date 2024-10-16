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

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	books, err := s.ValidateOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	o := &pb.Order{
		OrderId:    "order_1",
		Books:      books,
		CustomerId: p.CustomerId,
		Status:     "pending",
	}
	return o, nil

}

func (s *service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Book, error) {
	if len(req.Books) == 0 {
		return nil, shared.ErrNoBooks
	}

	mergedBooks := mergeBooksQuantities(req.Books)
	log.Println("Merged books:", mergedBooks)

	/// validate with the catalog service

	/// Temprorary
	var booksWithPrices []*pb.Book
	for _, b := range mergedBooks {
		booksWithPrices = append(booksWithPrices, &pb.Book{
			BookId:   b.BookId,
			PriceId:  "price_1QAUrkHHAM3KUbolwwfAqEjh",
			Quantity: b.Quantity,
		})
	}
	return booksWithPrices, nil
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
