package discovery

import (
	"context"
	"log"
	"math/rand"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ServiceConnection(ctx context.Context, serviceName string, registry DiscoveryRegistry) (*grpc.ClientConn, error) {
	addrs, err := registry.DiscoverService(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	log.Printf("Discovered %d instances of %s", len(addrs), serviceName)

	// Randomly select an instance
	return grpc.NewClient(
		addrs[rand.Intn(len(addrs))],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// Add OpenTelemetry interceptors
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),   // jaeger
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()), // jaeger
	)
}
