package gateway

import (
	"context"
	"log"

	"github.com/karokojnr/bookmesh-shared/discovery"
	pb "github.com/karokojnr/bookmesh-shared/proto"
)

type gateway struct {
	registry discovery.DiscoveryRegistry
}

func NewGrpcGateway(registry discovery.DiscoveryRegistry) *gateway {
	return &gateway{registry}
}

func (g *gateway) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {

	conn, err := discovery.ServiceConnection(context.Background(), "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	c := pb.NewOrderServiceClient(conn)

	return c.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerId: p.CustomerId,
		Books:      p.Books,
	})
}

func (g *gateway) GetOrder(ctx context.Context, orderId, customerId string) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	c := pb.NewOrderServiceClient(conn)

	return c.GetOrder(ctx, &pb.GetOrderRequest{
		OrderId:    orderId,
		CustomerId: customerId,
	})
}
