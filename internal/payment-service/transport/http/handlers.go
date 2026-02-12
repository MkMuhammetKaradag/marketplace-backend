// internal/payment-service/transport/http/handlers.go
package http

import (
	"fmt"
	"marketplace/internal/payment-service/domain"
	"marketplace/internal/payment-service/transport/http/controller"
	"marketplace/internal/payment-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Webhook *webhookHandlers
}

type webhookHandlers struct {
	Stripe *controller.StripeWebhookController
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: "hello payment service",
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}

func (h *Handlers) CreatePaymentSession(c *fiber.Ctx) error {
	fmt.Println("CreatePaymentSession")
	// paymentSessionRequest := domain.CreatePaymentSessionRequest{
	// 	OrderID:   uuid.New(),
	// 	UserID:    uuid.New(),
	// 	Amount:    10,
	// 	UserEmail: "test@mail.com",
	// }
	// fmt.Println("paymentSessionRequest", paymentSessionRequest)

	// paymentURL, err := h.stripeService.CreatePaymentSession(paymentSessionRequest)
	// if err != nil {
	// 	fmt.Println("error creating payment session", err)
	// 	return err
	// }

	return c.JSON(fiber.Map{
		"payment_url": "paymentURL",
	})
}

func NewHandlers(
	repo domain.PaymentRepository,
	stripeService domain.StripeService,
	messaging domain.Messaging,
) *Handlers {

	webhookUC := usecase.NewStripeWebhookUseCase(stripeService, messaging)

	return &Handlers{

		Webhook: &webhookHandlers{
			Stripe: controller.NewStripeWebhookController(webhookUC),
		},
	}
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
