package main

import (
	pb "github.com/karokojnr/bookmesh-shared/api"
)

type CreateOrderRequest struct {
	Order         *pb.Order `"json": order`
	RedirectToURL string    `"json": redirectToURL`
}
