package handler

import "net/http"

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerId}/orders", h.CreateOrder)
}

func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// get customer id from path
	// get order from request body
	// create order
	// return order
}
