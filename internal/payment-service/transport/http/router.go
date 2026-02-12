// internal/payment-service/transport/http/router.go
package http

import (
	"marketplace/internal/payment-service/handler"
	"marketplace/internal/payment-service/transport/http/controller"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	handlers *Handlers
}

func NewRouter(handlers *Handlers) *Router {
	return &Router{handlers: handlers}
}

func (r *Router) Register(app *fiber.App) {
	h := r.handlers

	// Public Routes
	app.Get("/hello", h.Hello)
	webhookController := r.handlers.Webhook.Stripe

	webhook := app.Group("/")
	{
		webhook.Post("/payment/webhook", handler.HandleWithFiber[controller.StripeWebhookRequest, controller.StripeWebhookResponse](webhookController))
	}
	app.Post("/create-payment-session", r.handlers.CreatePaymentSession)

}
