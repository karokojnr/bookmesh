package main

import (
	"context"
	"log"
	"net"
	"time"

	shared "github.com/karokojnr/bookmesh-shared"
	"github.com/karokojnr/bookmesh-shared/broker"
	"github.com/karokojnr/bookmesh-shared/discovery"
	"github.com/karokojnr/bookmesh-shared/discovery/consul"
	tracer "github.com/karokojnr/bookmesh-shared/tracer"

	"google.golang.org/grpc"
)

var (
	svcName    = "orders"
	consulAddr = shared.EnvString("CONSUL_ADDR", "localhost:8500")
	grpcAddr   = shared.EnvString("GRPC_ADDR", "localhost:2000")
	amqpUser   = shared.EnvString("RABBITMQ_USER", "guest")
	amqpPass   = shared.EnvString("RABBITMQ_PASS", "guest")
	amqpHost   = shared.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort   = shared.EnvString("RABBITMQ_PORT", "5672")
	jaegerAddr = shared.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {
	/// Tracer
	if err := tracer.SetGlobalTracer(context.TODO(), svcName, jaegerAddr); err != nil {
		log.Fatalf("Failed to set global tracer: %v", err)
	}

	/// Service Discovery Registry
	registry, err := consul.NewConsulDiscoveryRegistry(consulAddr, svcName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceID(svcName)
	if err := registry.RegisterService(ctx, instanceId, svcName, grpcAddr); err != nil {
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

	/// Broker
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	grpcServer := grpc.NewServer()
	conn, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer conn.Close()

	store := NewStorage()
	svc := NewService(store)

	/// Use decorator pattern to add middleware to the service
	svcWithTelemetryMiddleware := NewTelemetryMiddleware(svc)
	NewGrpcHandler(grpcServer, svcWithTelemetryMiddleware, ch)

	/// RabbitMQ consumer
	amqpConsumer := NewConsumer(svcWithTelemetryMiddleware)
	go amqpConsumer.Listen(ch)

	log.Println("Starting grpc server on", grpcAddr)
	if err := grpcServer.Serve(conn); err != nil {
		log.Fatalf(err.Error())
	}

}
