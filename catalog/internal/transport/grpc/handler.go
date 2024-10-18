package transport

import (
	"context"

	"github.com/karokojnr/bookmesh-catalog/internal/types"
	pb "github.com/karokojnr/bookmesh-shared/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type CatalogGrpcHandler struct {
	pb.UnimplementedCatalogServiceServer
	service types.CatalogService
	channel *amqp.Channel
}

func NewCatalogGrpcHandler(server *grpc.Server, channel *amqp.Channel, catalogService types.CatalogService) {
	h := &CatalogGrpcHandler{
		service: catalogService,
		channel: channel,
	}
	pb.RegisterCatalogServiceServer(server, h)
}

func (h *CatalogGrpcHandler) CheckIfBookIsInCatalog(ctx context.Context, req *pb.CheckIfBookIsInCatalogRequest) (*pb.CheckIfBookIsInCatalogResponse, error) {
	inCatalog, books, err := h.service.CheckIfBookIsInCatalog(ctx, req.Books)
	if err != nil {
		return nil, err
	}

	return &pb.CheckIfBookIsInCatalogResponse{
		IsInCatalog: inCatalog,
		Books:       books,
	}, nil
}

func (h *CatalogGrpcHandler) GetBooks(ctx context.Context, req *pb.GetBooksRequest) (*pb.GetBooksResponse, error) {
	books, err := h.service.GetBooks(ctx, req.BookIds)
	if err != nil {
		return nil, err
	}

	return &pb.GetBooksResponse{
		Books: books,
	}, nil
}
