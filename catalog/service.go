package main

import (
	"context"

	pb "github.com/karokojnr/bookmesh-shared/api"
)

type Service struct {
	store CatalogStore
}

func NewService(store CatalogStore) *Service {
	return &Service{store: store}
}

func (s *Service) CheckIfBookIsInCatalog(ctx context.Context, books []*pb.BookWithQuantity) (bool, []*pb.Book, error) {
	bookIds := make([]string, 0)
	for _, book := range books {
		bookIds = append(bookIds, book.BookId)
	}

	catalogBooks, err := s.store.GetBooks(ctx, bookIds)
	if err != nil {
		return false, nil, err
	}

	/// Check if all books are in the catalog
	for _, catalogBook := range catalogBooks {
		for _, book := range books {
			if catalogBook.BookId == book.BookId && catalogBook.Quantity < book.Quantity {
				return false, nil, nil
			}
		}
	}

	/// create books with prices from catalog
	res := make([]*pb.Book, 0)
	for _, catalogBook := range catalogBooks {
		for _, book := range books {
			if catalogBook.BookId == book.BookId {
				res = append(res, &pb.Book{
					BookId:   catalogBook.BookId,
					Title:    catalogBook.Title,
					Authors:  catalogBook.Authors,
					PriceId:  catalogBook.PriceId,
					Quantity: book.Quantity,
				})
			}
		}
	}

	return true, res, nil
}

func (s *Service) GetBooks(ctx context.Context, ids []string) ([]*pb.Book, error) {
	return s.store.GetBooks(ctx, ids)
}
