package types

import (
	"context"

	pb "github.com/karokojnr/bookmesh-shared/proto"
)

type CatalogService interface {
	CheckIfBookIsInCatalog(context.Context, []*pb.BookWithQuantity) (bool, []*pb.Book, error)
	GetBooks(ctx context.Context, ids []string) ([]*pb.Book, error)
}

type CatalogStore interface {
	GetBook(ctx context.Context, id string) (*pb.Book, error)
	GetBooks(ctx context.Context, ids []string) ([]*pb.Book, error)
}
