package domain

import "github.com/google/uuid"

type CreatePaymentSessionRequest struct {
	OrderID   uuid.UUID `json:"order_id"`
	UserID    uuid.UUID `json:"user_id"`
	Amount    float64   `json:"amount"`
	UserEmail string    `json:"user_email"`
	UserName  string    `json:"user_name"`
}

type CreatePaymentSessionResponse struct {
	PaymentURL string `json:"payment_url"`
}

type PaymentCompletedEvent struct {
	OrderID uuid.UUID `json:"order_id"`
	UserID  uuid.UUID `json:"user_id"`
	Status  string    `json:"status"`
}

type StripeService interface {
	CreatePaymentSession(req CreatePaymentSessionRequest) (string, error)
	GetWebhookSecret() string
}
