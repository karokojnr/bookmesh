package discovery

import (
	"context"
	"fmt"
	"time"

	"math/rand"
)

type DiscoveryRegistry interface {
	RegisterService(ctx context.Context, instanceId, serverName, hostPort string) error
	UnregisterService(ctx context.Context, instanceId, serviceName string) error
	DiscoverService(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(ctx context.Context, serviceName string) error
}

func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
