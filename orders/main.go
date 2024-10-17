package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/karokojnr/bookmesh-orders/gateway"
	shared "github.com/karokojnr/bookmesh-shared"
	"github.com/karokojnr/bookmesh-shared/broker"
	"github.com/karokojnr/bookmesh-shared/discovery"
	"github.com/karokojnr/bookmesh-shared/discovery/consul"
	tracer "github.com/karokojnr/bookmesh-shared/tracer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

var (
	svcName     = "orders"
	consulAddr  = shared.EnvString("CONSUL_ADDR", "localhost:8500")
	grpcAddr    = shared.EnvString("GRPC_ADDR", "localhost:2000")
	amqpUser    = shared.EnvString("RABBITMQ_USER", "guest")
	amqpPass    = shared.EnvString("RABBITMQ_PASS", "guest")
	amqpHost    = shared.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort    = shared.EnvString("RABBITMQ_PORT", "5672")
	jaegerAddr  = shared.EnvString("JAEGER_ADDR", "localhost:4318")
	mongoDbUser = shared.EnvString("MONGO_DB_USER", "root")
	mongoDbPass = shared.EnvString("MONGO_DB_PASS", "example")
	mongoDbAddr = shared.EnvString("MONGO_DB_HOST", "localhost:27017")
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	/// Tracer
	if err := tracer.SetGlobalTracer(context.TODO(), svcName, jaegerAddr); err != nil {
		logger.Fatal("Failed to set global tracer ", zap.Error(err))
	}

	/// Service Discovery Registry
	registry, err := consul.NewConsulDiscoveryRegistry(consulAddr, svcName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceID(svcName)
	if err := registry.RegisterService(ctx, instanceId, svcName, grpcAddr); err != nil {
		logger.Fatal("Failed to register service ", zap.Error(err))
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceId, svcName); err != nil {
				logger.Fatal("Health check failed ", zap.Error(err))
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

	/// Mongo db connection
	// mongo db conn
	uri := fmt.Sprintf("mongodb://%s:%s@%s", mongoDbUser, mongoDbPass, mongoDbAddr)
	mongoClient, err := connectToMongoDb(uri)
	if err != nil {
		logger.Fatal("failed to connect to mongo db", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	conn, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatal("Failed to listen ", zap.Error(err))
	}
	defer conn.Close()

	/// Gateway
	gateway := gateway.NewGateway(registry)

	store := NewStorage(mongoClient)
	svc := NewService(store, gateway)

	/// Use decorator pattern to add middleware to the service
	svcWithTelemetryMiddleware := NewTelemetryMiddleware(svc)
	svcWithLoggingMiddleware := NewLoggingMiddleware(svcWithTelemetryMiddleware)
	NewGrpcHandler(grpcServer, svcWithLoggingMiddleware, ch)

	/// RabbitMQ consumer
	amqpConsumer := NewConsumer(svcWithLoggingMiddleware)
	go amqpConsumer.Listen(ch)

	logger.Info("Starting grpc server on", zap.String("address", grpcAddr))
	if err := grpcServer.Serve(conn); err != nil {
		logger.Fatal("Failed to start grpc server ", zap.Error(err))
	}

}

func connectToMongoDb(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	return client, err
}
