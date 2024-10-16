package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/karokojnr/bookmesh-gateway/gateway"
	shared "github.com/karokojnr/bookmesh-shared"
	"github.com/karokojnr/bookmesh-shared/discovery"
	"github.com/karokojnr/bookmesh-shared/discovery/consul"
	tracer "github.com/karokojnr/bookmesh-shared/tracer"
)

var (
	httpAddr   = shared.EnvString("HTTP_ADDR", ":8080")
	consulAddr = shared.EnvString("CONSUL_ADDR", "localhost:8500")
	svcName    = "gateway"
	jaegerAddr = shared.EnvString("JAEGER_ADDR", "http://localhost:4318")
)

func main() {
	/// Tracer
	if err := tracer.SetGlobalTracer(context.TODO(), svcName, jaegerAddr); err != nil {
		log.Fatalf("Failed to set global tracer: %v", err)
	}

	/// Register consul
	registry, err := consul.NewConsulDiscoveryRegistry(consulAddr, svcName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceID(svcName)
	if err := registry.RegisterService(ctx, instanceId, svcName, httpAddr); err != nil {
		log.Fatalf("Failed to register service: %v", err)
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
	///

	mux := http.NewServeMux()

	ordersGateway := gateway.NewGrpcGateway(registry)
	h := NewHttpHandler(ordersGateway)
	h.RegisterRoutes(mux)

	log.Println("Starting http server on", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server", err)
	}
}
