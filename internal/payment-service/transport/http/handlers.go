// internal/payment-service/transport/http/handlers.go
package http

import (
	"encoding/json"
	"fmt"
	"marketplace/internal/payment-service/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type Handlers struct {
	stripeService domain.StripeService
}

func NewHandlers(stripeService domain.StripeService) *Handlers {
	return &Handlers{
		stripeService: stripeService,
	}
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
	paymentSessionRequest := domain.CreatePaymentSessionRequest{
		OrderID:   uuid.New(),
		UserID:    uuid.New(),
		Amount:    10,
		UserEmail: "test@mail.com",
	}
	fmt.Println("paymentSessionRequest", paymentSessionRequest)

	paymentURL, err := h.stripeService.CreatePaymentSession(paymentSessionRequest)
	if err != nil {
		fmt.Println("error creating payment session", err)
		return err
	}

	return c.JSON(fiber.Map{
		"payment_url": paymentURL,
	})
}

func (h *Handlers) StripeWebhook(c *fiber.Ctx) error {
	payload := c.Body()
	sigHeader := c.Get("Stripe-Signature")
	endpointSecret := h.stripeService.GetWebhookSecret()

	// 1. Gelen isteğin gerçekten Stripe'tan geldiğini doğrula
	event, err := webhook.ConstructEventWithOptions(
		payload,
		sigHeader,
		endpointSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true, // Hatanın çözümü tam olarak burası
		},
	)
	if err != nil {
		fmt.Printf("⚠️ Webhook doğrulama hatası: %v\n", err)
		return c.Status(400).SendString("Geçersiz imza")
	}

	// 2. Sadece "Ödeme Başarılı" (checkout.session.completed) olayıyla ilgileniyoruz
	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			return c.Status(400).SendString("Data parse hatası")
		}

		// Stripe Session oluştururken içine koyduğumuz Metadata'yı geri alıyoruz
		orderID := session.Metadata["order_id"]
		userID := session.Metadata["user_id"]

		fmt.Printf("✅ Ödeme Başarılı! Order ID: %s, User ID: %s\n", orderID, userID)

		// 3. ŞİMDİ SIRADAKİ HAMLE: KAFKA'YA MESAJ ATMAK
		// Buraya birazdan Kafka Producer kodunu bağlayacağız
		// h.kafkaProducer.PublishPaymentSuccess(orderID, userID)
	}

	return c.SendStatus(200)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
