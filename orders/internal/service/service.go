package service

import (
	"context"

	"github.com/karokojnr/bookmesh-orders/internal/catalog_gateway"
	"github.com/karokojnr/bookmesh-orders/internal/types"
	shared "github.com/karokojnr/bookmesh-shared"
	pb "github.com/karokojnr/bookmesh-shared/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrdersStore interface {
	Create(context.Context, types.Order) (primitive.ObjectID, error)
	Get(ctx context.Context, id, customerID string) (*types.Order, error)
	Update(ctx context.Context, id string, o *pb.Order) error
}

type service struct {
	store   OrdersStore
	gateway catalog_gateway.CatalogGateway
}

func NewService(store OrdersStore, gateway catalog_gateway.CatalogGateway) *service {
	return &service{
		store:   store,
		gateway: gateway,
	}
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, books []*pb.Book) (*pb.Order, error) {

	id, err := s.store.Create(ctx, types.Order{
		CustomerId:  p.CustomerId,
		Status:      "pending",
		Books:       books,
		PaymentLink: "",
	})
	if err != nil {
		return nil, err
	}

	o := &pb.Order{
		OrderId:    id.Hex(),
		CustomerId: p.CustomerId,
		Status:     "pending",
		Books:      books,
	}

	return o, nil

}

func (s *service) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	o, err := s.store.Get(ctx, req.OrderId, req.CustomerId)
	if err != nil {
		return nil, err
	}

	return o.ToProto(), nil
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

	//  validate with the catalog service
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
