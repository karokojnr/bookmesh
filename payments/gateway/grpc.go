package gateway

import (
	"context"
	"log"

	pb "github.com/karokojnr/bookmesh-shared/api"
	"github.com/karokojnr/bookmesh-shared/discovery"
)

type gateway struct {
	registry discovery.DiscoveryRegistry
}

func NewGateway(registry discovery.DiscoveryRegistry) *gateway {
	return &gateway{registry}
}

func (g *gateway) UpdateOrderWithPaymentLink(ctx context.Context, orderId, link string) error {
	conn, err := discovery.ServiceConnection(context.Background(), "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	ordersClient := pb.NewOrderServiceClient(conn)

	_, err = ordersClient.UpdateOrder(ctx, &pb.Order{
		OrderId:          orderId,
		Status:      "waiting_payment",
		PaymentLink: link,
	})
	return err
}
