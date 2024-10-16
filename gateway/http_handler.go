package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/karokojnr/bookmesh-gateway/gateway"
	shared "github.com/karokojnr/bookmesh-shared"
	pb "github.com/karokojnr/bookmesh-shared/api"
	"go.opentelemetry.io/otel"
	otelCodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type httpHandler struct {
	gateway gateway.OrdersGateway
}

func NewHttpHandler(gateway gateway.OrdersGateway) *httpHandler {
	return &httpHandler{gateway}
}

func (h *httpHandler) RegisterRoutes(mux *http.ServeMux) {
	/// static file server
	mux.Handle("/", http.FileServer(http.Dir("public")))

	mux.HandleFunc("POST /api/customers/{customerId}/orders", h.createOrder)
	mux.HandleFunc("GET /api/customers/{customerId}/orders/{orderId}", h.getOrder)

}

func (h *httpHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customerId")

	var books []*pb.BookWithQuantity

	if err := shared.ReadJSON(r, &books); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Trace
	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	if err := validateBooks(books); err != nil {
		shared.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.gateway.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerId: customerId,
		Books:      books,
	})

	/// grpc error handling
	errStatus := status.Convert(err)

	if errStatus != nil {
		span.SetStatus(otelCodes.Error, err.Error())
		if errStatus.Code() == codes.InvalidArgument {
			shared.WriteError(w, http.StatusBadRequest, errStatus.Message())
			return
		}

		shared.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res := &CreateOrderRequest{
		Order:         o,
		RedirectToURL: fmt.Sprintf("http://localhost:8080/success.html?customerId=%s&orderId=%s", o.CustomerId, o.OrderId),
	}

	shared.WriteJSON(w, http.StatusCreated, res)
}

func (h *httpHandler) getOrder(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customerId")
	orderId := r.PathValue("orderId")

	// Trace
	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	o, err := h.gateway.GetOrder(ctx, orderId, customerId)

	/// grpc error handling
	errStatus := status.Convert(err)

	if errStatus != nil {
		span.SetStatus(otelCodes.Error, err.Error())

		if errStatus.Code() != codes.InvalidArgument {
			shared.WriteError(w, http.StatusBadRequest, errStatus.Message())
			return
		}

		shared.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	shared.WriteJSON(w, http.StatusOK, o)

}

func validateBooks(books []*pb.BookWithQuantity) error {
	if len(books) == 0 {
		return shared.ErrNoBooks
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
