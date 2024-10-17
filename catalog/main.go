package main

import (
	"context"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"
	shared "github.com/karokojnr/bookmesh-shared"
	"github.com/karokojnr/bookmesh-shared/broker"
	"github.com/karokojnr/bookmesh-shared/discovery"
	"github.com/karokojnr/bookmesh-shared/discovery/consul"
	tracer "github.com/karokojnr/bookmesh-shared/tracer"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceName = "catalog"
	grpcAddr    = shared.EnvString("GRPC_ADDR", "localhost:2002")
	consulAddr  = shared.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser    = shared.EnvString("RABBITMQ_USER", "guest")
	amqpPass    = shared.EnvString("RABBITMQ_PASS", "guest")
	amqpHost    = shared.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort    = shared.EnvString("RABBITMQ_PORT", "5672")
	jaegerAddr  = shared.EnvString("JAEGER_ADDR", "localhost:4318")
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	if err := tracer.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr); err != nil {
		logger.Fatal("could set global tracer", zap.Error(err))
	}

	registry, err := consul.NewConsulDiscoveryRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.RegisterService(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				logger.Error("Failed to health check", zap.Error(err))
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.UnregisterService(ctx, instanceID, serviceName)

	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	defer l.Close()

	store := NewStore()
	svc := NewService(store)
	svcWithTelemetry := NewTelemetryMiddleware(svc)

	NewCatalogGrpcHandler(grpcServer, ch, svcWithTelemetry)

	consumer := NewConsumer()
	go consumer.Listen(ch)

	logger.Info("Starting gRPC server", zap.String("port", grpcAddr))

	if err := grpcServer.Serve(l); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
