package main

import (
	"log"
	"net"

	shared "github.com/karokojnr/bookmesh-shared"
	"google.golang.org/grpc"
)

var (
	grpcAddr = shared.EnvString("GRPC_ADDR", "localhost:8081")
)

func main() {

	grpcServer := grpc.NewServer()
	conn, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer conn.Close()

	store := NewStorage()
	svc := NewService(store)
	NewGrpcHandler(grpcServer)

	log.Println("Starting grpc server on", grpcAddr)
	svc.CreateOrder()

	if err := grpcServer.Serve(conn); err != nil {
		log.Fatalf(err.Error())
	}

}
