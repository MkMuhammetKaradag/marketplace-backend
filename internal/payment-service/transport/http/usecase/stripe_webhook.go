package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace/internal/payment-service/domain"
	eventsProto "marketplace/pkg/proto/events"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type StripeWebhookUseCase interface {
	Execute(ctx context.Context, payload []byte, sigHeader string) error
}

type stripeWebhookUseCase struct {
	stripeService domain.StripeService
	messaging     domain.Messaging
}

func NewStripeWebhookUseCase(stripeService domain.StripeService, messaging domain.Messaging) StripeWebhookUseCase {
	return &stripeWebhookUseCase{
		stripeService: stripeService,
		messaging:     messaging,
	}
}

func (u *stripeWebhookUseCase) Execute(ctx context.Context, payload []byte, sigHeader string) error {

	endpointSecret := u.stripeService.GetWebhookSecret()

	event, err := webhook.ConstructEventWithOptions(
		payload,
		sigHeader,
		endpointSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		return err
	}

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			return err
		}

		orderID := session.Metadata["order_id"]
		userID := session.Metadata["user_id"]

		fmt.Printf("✅ Ödeme Başarılı! Order ID: %s, User ID: %s\n", orderID, userID)

		msg := &eventsProto.Message{
			Type:        eventsProto.MessageType_PAYMENT_SUCCESSFUL,
			FromService: eventsProto.ServiceType_PAYMENT_SERVICE,
			RetryCount:  5,
			ToServices:  []eventsProto.ServiceType{eventsProto.ServiceType_ORDER_SERVICE, eventsProto.ServiceType_BASKET_SERVICE, eventsProto.ServiceType_PRODUCT_SERVICE},
			Payload: &eventsProto.Message_PaymentSuccessfulData{PaymentSuccessfulData: &eventsProto.PaymentSuccessfulData{
				OrderId:         orderID,
				UserId:          userID,
				Amount:          float64(session.AmountTotal),
				StripeSessionId: session.ID,
			}},
		}

		u.messaging.PublishMessage(ctx, msg)
	}

	return nil

}
