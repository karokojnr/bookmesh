package consul

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
)

type ConsulDiscoveryRegistry struct {
	client *consul.Client
}

func NewConsulDiscoveryRegistry(addr, serviceName string) (*ConsulDiscoveryRegistry, error) {
	config := consul.DefaultConfig()
	config.Address = addr

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConsulDiscoveryRegistry{client: client}, nil
}

func (r ConsulDiscoveryRegistry) RegisterService(ctx context.Context, instanceID, serviceName, hostPort string) error {
	host, portStr, found := strings.Cut(hostPort, ":")
	if !found {
		return errors.New("invalid host:port format. Eg: localhost:8081")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID:      instanceID,
		Address: host,
		Port:    port,
		Name:    serviceName,
		Check: &consul.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  true,
			TTL:                            "5s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}

func (r ConsulDiscoveryRegistry) UnregisterService(ctx context.Context, instanceID string, serviceName string) error {
	log.Printf("Deregistering service %s", instanceID)
	return r.client.Agent().CheckDeregister(instanceID)
}

func (r ConsulDiscoveryRegistry) DiscoverService(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	var instances []string
	for _, entry := range entries {
		instances = append(instances, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
	}

	return instances, nil
}

func (r ConsulDiscoveryRegistry) HealthCheck(instanceID string, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", consul.HealthPassing)
}
