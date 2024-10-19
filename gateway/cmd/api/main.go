package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/karokojnr/bookmesh-gateway/internal/gateway"
	transport "github.com/karokojnr/bookmesh-gateway/internal/transport/http"
	shared "github.com/karokojnr/bookmesh-shared"
	"github.com/karokojnr/bookmesh-shared/discovery"
	"github.com/karokojnr/bookmesh-shared/discovery/consul"
	tracer "github.com/karokojnr/bookmesh-shared/tracer"
)

var (
	svcName    = "gateway"
	httpAddr   = shared.EnvString("HTTP_ADDR", ":8080")
	consulAddr = shared.EnvString("CONSUL_ADDR", "localhost:8500")
	jaegerAddr = shared.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {
	// Tracer
	if err := tracer.SetGlobalTracer(context.TODO(), svcName, jaegerAddr); err != nil {
		log.Fatalf("failed to set global tracer in gateway: %v", err)
	}

	// Register consul
	registry, err := consul.NewConsulDiscoveryRegistry(consulAddr, svcName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceID(svcName)
	if err := registry.RegisterService(ctx, instanceId, svcName, httpAddr); err != nil {
		log.Fatalf("failed to register gateway service: %v", err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceId, svcName); err != nil {
				log.Fatalf("Health check failed: %v", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	defer registry.UnregisterService(ctx, instanceId, svcName)

	mux := http.NewServeMux()

	ordersGateway := gateway.NewGrpcGateway(registry)
	h := transport.NewHttpHandler(ordersGateway)
	h.RegisterRoutes(mux)

	log.Println("gateway starting http server on: ", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("gateway failed to start http server", err)
	}
}
