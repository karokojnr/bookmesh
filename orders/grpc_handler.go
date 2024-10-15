package main

import (
	"context"
	"log"

	pb "github.com/karokojnr/bookmesh-shared/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	svc OrdersService
}

func NewGrpcHandler(grpcServer *grpc.Server, svc OrdersService) {
	handler := &grpcHandler{
		svc: svc,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	/// show gRPC error handling
	/// return nil, fmt.Errorf("not implemented")

	log.Printf("Creating order for customer %v", req)
	o := &pb.Order{
		OrderId: "1",
	}
	return o, nil
}
