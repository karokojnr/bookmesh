package main

import (
	"context"
	"log"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"

	stripeprocessor "github.com/karokojnr/bookmesh-payments/payment_processor/stripe_processor"
	shared "github.com/karokojnr/bookmesh-shared"
	"github.com/karokojnr/bookmesh-shared/broker"
	"github.com/karokojnr/bookmesh-shared/discovery"
	"github.com/karokojnr/bookmesh-shared/discovery/consul"
	"github.com/stripe/stripe-go/v78"
	"google.golang.org/grpc"
)

var (
	svcName    = "payments"
	consulAddr = shared.EnvString("CONSUL_ADDR", "localhost:8500")
	grpcAddr   = shared.EnvString("GRPC_ADDR", "localhost:8082")
	amqpUser   = shared.EnvString("RABBITMQ_USER", "guest")
	amqpPass   = shared.EnvString("RABBITMQ_PASS", "guest")
	amqpHost   = shared.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort   = shared.EnvString("RABBITMQ_PORT", "5672")
	stripeKey  = shared.EnvString("STRIPE_KEY", "")
)

func main() {
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
	log.Println("Stripe key: ", stripeKey)
	stripeProcessor := stripeprocessor.NewStripe()

	/// Broker
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	svc := NewService(stripeProcessor)

	amqpConsumer := NewConsumer(svc)
	go amqpConsumer.Listen(ch)

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
