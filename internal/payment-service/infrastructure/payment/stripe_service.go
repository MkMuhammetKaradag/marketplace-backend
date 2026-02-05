package payment

import (
	"marketplace/internal/payment-service/domain"
	"time"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/checkout/session"
)

type StripeService struct {
	secretKey     string
	webhookSecret string
}

func NewStripeService(key string, webhookSecret string) *StripeService {
	stripe.Key = key
	return &StripeService{secretKey: key, webhookSecret: webhookSecret}
}

func (s *StripeService) CreatePaymentSession(req domain.CreatePaymentSessionRequest) (string, error) {
	expiresAt := time.Now().Add(30 * time.Minute).Unix()
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		CustomerEmail:      stripe.String(req.UserEmail),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Order #" + req.OrderID.String()),
					},
					UnitAmount: stripe.Int64(int64(req.Amount * 100)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("http://localhost:3000/success?order_id=" + req.OrderID.String()),
		CancelURL:  stripe.String("http://localhost:3000/cancel"),
		Metadata: map[string]string{
			"order_id":   req.OrderID.String(),
			"user_id":    req.UserID.String(),
			"user_name":  req.UserName,
			"user_email": req.UserEmail,
		},
		ExpiresAt: stripe.Int64(expiresAt),
	}

	sess, err := session.New(params)
	if err != nil {
		return "", err
	}

	return sess.URL, nil
}

func (s *StripeService) GetWebhookSecret() string {
	return s.webhookSecret
}
