package gateway

import "context"

type OrdersGateway interface {
	UpdateOrderWithPaymentLink(ctx context.Context, orderId, link string) error
}
