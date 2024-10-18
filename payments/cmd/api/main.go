package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/karokojnr/bookmesh-payments/gateway"
	"github.com/karokojnr/bookmesh-payments/internal/consumer"
	stripeprocessor "github.com/karokojnr/bookmesh-payments/internal/payment_processor/stripe_processor"
	"github.com/karokojnr/bookmesh-payments/internal/service"
	transport "github.com/karokojnr/bookmesh-payments/internal/transport/http"
	middleware "github.com/karokojnr/bookmesh-payments/internal/transport/middlewares"
	shared "github.com/karokojnr/bookmesh-shared"
	"github.com/karokojnr/bookmesh-shared/broker"
	"github.com/karokojnr/bookmesh-shared/discovery"
	"github.com/karokojnr/bookmesh-shared/discovery/consul"
	tracer "github.com/karokojnr/bookmesh-shared/tracer"
	"github.com/stripe/stripe-go/v78"
	"google.golang.org/grpc"
)

var (
	svcName    = "payments"
	consulAddr = shared.EnvString("CONSUL_ADDR", "localhost:8500")
	grpcAddr   = shared.EnvString("GRPC_ADDR", "localhost:2001")
	httpAddr   = shared.EnvString("HTTP_ADDR", "localhost:8081")
	amqpUser   = shared.EnvString("RABBITMQ_USER", "guest")
	amqpPass   = shared.EnvString("RABBITMQ_PASS", "guest")
	amqpHost   = shared.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort   = shared.EnvString("RABBITMQ_PORT", "5672")
	stripeKey  = shared.EnvString("STRIPE_KEY", "")
	jaegerAddr = shared.EnvString("JAEGER_ADDR", "localhost:4318")
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

	/// Stripe
	stripe.Key = stripeKey
	stripeProcessor := stripeprocessor.NewStripe()

	/// Broker
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	/// Gateway
	gateway := gateway.NewGateway(registry)
	///
	/// Use decorator pattern to add middleware to the service

	svc := service.NewService(stripeProcessor, gateway)
	svcWithTelemetryMiddleware := middleware.NewTelemetryMiddleware(svc)

	amqpConsumer := consumer.NewConsumer(svcWithTelemetryMiddleware)
	go amqpConsumer.Listen(ch)

	/// Http server
	mux := http.NewServeMux()
	httpServer := transport.NewHttpPaymentHandler(ch)
	httpServer.RegisterRoutes(mux)
	go func() {
		log.Println("Starting http server on ", httpAddr)
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			log.Fatalf("Failed to start http server: %v", err)
		}
	}()

	/// gRPC server
	grpcServer := grpc.NewServer()
	conn, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer conn.Close()

	log.Println("Starting grpc server on", grpcAddr)

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatalf(err.Error())
	}

}
