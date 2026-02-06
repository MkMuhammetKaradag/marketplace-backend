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
	var session stripe.CheckoutSession
	if event.Type == "checkout.session.completed" ||
		event.Type == "checkout.session.expired" ||
		event.Type == "checkout.session.async_payment_failed" {

		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			return err
		}
	}

	orderID := session.Metadata["order_id"]
	userID := session.Metadata["user_id"]
	userName := session.Metadata["user_name"]
	userEmail := session.Metadata["user_email"]
	amount := float64(session.AmountTotal)

	if orderID == "" || userID == "" {
		return fmt.Errorf("order id or user id is empty")
	}
	switch event.Type {
	case "checkout.session.completed":

		//return u.handleFailure(ctx, orderID, userID, string(event.Type))
		return u.handleSuccessful(ctx, orderID, userID, userName, userEmail, amount)

	case "checkout.session.expired", "checkout.session.async_payment_failed":
		return u.handleFailure(ctx, orderID, userID, userName, userEmail, string(event.Type))
	}

	return nil

}

func (u *stripeWebhookUseCase) handleFailure(ctx context.Context, orderID, userID, userName, userEmail, eventType string) error {
	fmt.Printf("❌ Ödeme Başarısız veya Süre Doldu! Order ID: %s\n", orderID)

	msg := &eventsProto.Message{
		Type:        eventsProto.MessageType_PAYMENT_FAILED,
		FromService: eventsProto.ServiceType_PAYMENT_SERVICE,
		ToServices:  []eventsProto.ServiceType{eventsProto.ServiceType_ORDER_SERVICE, eventsProto.ServiceType_BASKET_SERVICE, eventsProto.ServiceType_PRODUCT_SERVICE},
		Payload: &eventsProto.Message_PaymentFailedData{PaymentFailedData: &eventsProto.PaymentFailedData{
			OrderId:      orderID,
			ErrorMessage: eventType,
			UserId:       userID,
			FailureCode:  "",
		}},
	}
	return u.messaging.PublishMessage(ctx, msg)
}

func (u *stripeWebhookUseCase) handleSuccessful(ctx context.Context, orderID, userID, userName, userEmail string, amount float64) error {
	fmt.Printf("✅ Ödeme Başarılı! Order ID: %s, User ID: %s\n", orderID, userID)

	msg := &eventsProto.Message{
		Type:        eventsProto.MessageType_PAYMENT_SUCCESSFUL,
		FromService: eventsProto.ServiceType_PAYMENT_SERVICE,
		RetryCount:  5,
		ToServices:  []eventsProto.ServiceType{eventsProto.ServiceType_ORDER_SERVICE, eventsProto.ServiceType_BASKET_SERVICE, eventsProto.ServiceType_PRODUCT_SERVICE, eventsProto.ServiceType_NOTIFICATION_SERVICE},
		Payload: &eventsProto.Message_PaymentSuccessfulData{PaymentSuccessfulData: &eventsProto.PaymentSuccessfulData{
			OrderId:         orderID,
			UserId:          userID,
			Amount:          amount,
			StripeSessionId: "session.ID",
		}},
	}

	return u.messaging.PublishMessage(ctx, msg)
}
