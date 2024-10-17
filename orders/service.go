package main

import (
	"context"

	"github.com/karokojnr/bookmesh-orders/gateway"
	shared "github.com/karokojnr/bookmesh-shared"
	pb "github.com/karokojnr/bookmesh-shared/api"
)

type service struct {
	store   OrdersStore
	gateway gateway.CatalogGateway
}

func NewService(store OrdersStore, gateway gateway.CatalogGateway) *service {
	return &service{
		store:   store,
		gateway: gateway,
	}
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, books []*pb.Book) (*pb.Order, error) {

	id, err := s.store.Create(ctx, p, books)
	if err != nil {
		return nil, err
	}
	o := &pb.Order{
		OrderId:    id,
		CustomerId: p.CustomerId,
		Books:      books,
	}

	return o, nil

}

func (s *service) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	return s.store.Get(ctx, req.OrderId, req.CustomerId)
}

func (s *service) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	err := s.store.Update(ctx, o.OrderId, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (s *service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Book, error) {
	if len(req.Books) == 0 {
		return nil, shared.ErrNoBooks
	}

	mergedBooks := mergeBooksQuantities(req.Books)

	///  validate with the catalog service
	isInCatalog, books, err := s.gateway.CheckIfBookIsInCatalog(ctx, req.CustomerId, mergedBooks)

	if err != nil {
		return nil, err
	}
	if !isInCatalog {
		return books, shared.ErrNoCatalog
	}

	return books, nil
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
