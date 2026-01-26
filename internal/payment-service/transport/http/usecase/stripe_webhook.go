package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace/internal/payment-service/domain"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type StripeWebhookUseCase interface {
	Execute(ctx context.Context, payload []byte, sigHeader string) error
}

type stripeWebhookUseCase struct {
	stripeService domain.StripeService
}

func NewStripeWebhookUseCase(stripeService domain.StripeService) StripeWebhookUseCase {
	return &stripeWebhookUseCase{
		stripeService: stripeService,
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
		// fmt.Println("body:", payload)
		// 3. ŞİMDİ SIRADAKİ HAMLE: KAFKA'YA MESAJ ATMAK
		// Buraya birazdan Kafka Producer kodunu bağlayacağız
		// h.kafkaProducer.PublishPaymentSuccess(orderID, userID)
	}

	return nil

}
