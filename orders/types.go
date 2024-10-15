package main

import "context"

type OrdersService interface {
	CreateOrder(context.Context) error
}

type OrdersStore interface {
	Create(context.Context) error
}
