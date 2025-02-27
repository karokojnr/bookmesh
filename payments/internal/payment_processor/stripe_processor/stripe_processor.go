package stripeprocessor

import (
	"fmt"
	"log"

	shared "github.com/karokojnr/bookmesh-shared"
	pb "github.com/karokojnr/bookmesh-shared/proto"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

var gatewayAddr = shared.EnvString("GATEWAY_HTTP_ADDR", "http://localhost:8080")

type StripeProcessor struct{}

func NewStripe() *StripeProcessor {
	return &StripeProcessor{}
}

func (s *StripeProcessor) CreatePaymentLink(o *pb.Order) (string, error) {
	log.Printf("Creating payment link for order: %v", o)
	gatewaySuccessURL := fmt.Sprintf("%s/success.html?customerId=%s&orderId=%s", gatewayAddr, o.CustomerId, o.OrderId)
	gatewayCancelURL := fmt.Sprintf("%s/cancel.html", gatewayAddr)

	books := []*stripe.CheckoutSessionLineItemParams{}
	for _, b := range o.Books {
		books = append(books, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(b.PriceId),
			Quantity: stripe.Int64(int64(b.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		LineItems:  books,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(gatewaySuccessURL),
		CancelURL:  stripe.String(gatewayCancelURL),
		Metadata: map[string]string{
			"orderId":    o.OrderId,
			"customerId": o.CustomerId,
		},
	}

	result, err := session.New(params)
	if err != nil {
		return "", nil
	}

	return result.URL, nil

}
