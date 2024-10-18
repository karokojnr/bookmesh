package catalog_gateway

import (
	"context"

	pb "github.com/karokojnr/bookmesh-shared/proto"
)

type CatalogGateway interface {
	CheckIfBookIsInCatalog(ctx context.Context, customerId string, books []*pb.BookWithQuantity) (bool, []*pb.Book, error)
}
