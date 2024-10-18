package types

import (
	pb "github.com/karokojnr/bookmesh-shared/proto"
)

type CreateOrderRequest struct {
	Order         *pb.Order `json:"order"`
	RedirectToUrl string    `json:"redirect_to_url"`
}
