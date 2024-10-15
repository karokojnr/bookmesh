package main

import (
	"errors"
	"net/http"

	shared "github.com/karokojnr/bookmesh-shared"
	pb "github.com/karokojnr/bookmesh-shared/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type httpHandler struct {
	client pb.OrderServiceClient
}

func NewHttpHandler(c pb.OrderServiceClient) *httpHandler {
	return &httpHandler{client: c}
}

func (h *httpHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerId}/orders", h.CreateOrder)
}

func (h *httpHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customerId")

	var books []*pb.BookWithQuantity

	if err := shared.ReadJSON(r, &books); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateBooks(books); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerId: customerId,
		Books:      books,
	})

	/// grpc error handling
	errStatus := status.Convert(err)

	if errStatus != nil {
		if errStatus.Code() == codes.InvalidArgument {
			shared.WriteError(w, http.StatusBadRequest, errStatus.Message())
			return
		}

		shared.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	shared.WriteJSON(w, http.StatusCreated, o)
}

func validateBooks(books []*pb.BookWithQuantity) error {
	if len(books) == 0 {
		return errors.New("books must not be empty")
	}

	for _, b := range books {
		if b.BookId == "" {
			return errors.New("book id is required")
		}

		if b.Quantity <= 0 {
			return errors.New("invalid quantity")
		}
	}

	return nil
}
