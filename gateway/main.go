package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	shared "github.com/karokojnr/bookmesh-shared"
	pb "github.com/karokojnr/bookmesh-shared/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr         = shared.EnvString("HTTP_ADDR", ":8080")
	orderServiceAddr = "localhost:8081"
)

func main() {
	/// grpc client
	conn, err := grpc.NewClient(orderServiceAddr, grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatal("Failed to dial order service", err)
	}
	defer conn.Close()
	log.Println("Connected to order service at", orderServiceAddr)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	h := NewHttpHandler(c)
	h.RegisterRoutes(mux)

	log.Println("Starting http server on", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server", err)
	}
}
