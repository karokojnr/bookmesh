package gateway

import (
	"context"
	"log"

	pb "github.com/karokojnr/bookmesh-shared/api"
	"github.com/karokojnr/bookmesh-shared/discovery"
)

type Gateway struct {
	registry discovery.DiscoveryRegistry
}

func NewGateway(registry discovery.DiscoveryRegistry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) CheckIfBookIsInCatalog(ctx context.Context, customerId string, books []*pb.BookWithQuantity) (bool, []*pb.Book, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "catalog", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()
	client := pb.NewCatalogServiceClient(conn)

	res, err := client.CheckIfBookIsInCatalog(ctx, &pb.CheckIfBookIsInCatalogRequest{
		Books: books,
	})

	return res.IsInCatalog, res.Books, err
}
