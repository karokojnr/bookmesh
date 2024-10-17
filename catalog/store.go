package main

import (
	"context"
	"fmt"

	pb "github.com/karokojnr/bookmesh-shared/api"
)

type Store struct {
	catalog map[string]*pb.Book
}

func NewStore() *Store {
	return &Store{
		catalog: map[string]*pb.Book{
			"1": {
				BookId:   "1",
				Title:    "The Alchemist",
				Authors:  []string{"Paulo Coelho"},
				PriceId:  "price_1QAnwNHHAM3KUbolVHLhtu87",
				Quantity: 4,
			},
			"2": {
				BookId:   "2",
				Title:    "The Richest Man in Babylon",
				Authors:  []string{"George S. Clason"},
				PriceId:  "price_1QAUrkHHAM3KUbolwwfAqEjh",
				Quantity: 3,
			},
		},
	}
}

func (s *Store) GetBook(ctx context.Context, id string) (*pb.Book, error) {
	for _, book := range s.catalog {
		if book.BookId == id {
			return book, nil
		}
	}
	return nil, fmt.Errorf("book not found")
}

func (s *Store) GetBooks(ctx context.Context, ids []string) ([]*pb.Book, error) {
	var books []*pb.Book
	for _, id := range ids {
		if book, ok := s.catalog[id]; ok {
			books = append(books, book)
		}
	}
	return books, nil
}
