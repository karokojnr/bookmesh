package gateway

import (
	"context"

	pb "github.com/karokojnr/bookmesh-shared/api"
)

type CatalogGateway interface {
	CheckIfBookIsInCatalog(ctx context.Context, customerId string, books []*pb.BookWithQuantity) (bool, []*pb.Book, error)
}
