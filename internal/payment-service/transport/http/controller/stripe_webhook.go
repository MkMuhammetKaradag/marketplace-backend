package controller

import (
	"marketplace/internal/payment-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type StripeWebhookRequest struct {
	Signature string `reqHeader:"Stripe-Signature"`
}

type StripeWebhookResponse struct {
	Message string `json:"message"`
}
type StripeWebhookController struct {
	usecase usecase.StripeWebhookUseCase
}

func NewStripeWebhookController(usecase usecase.StripeWebhookUseCase) *StripeWebhookController {
	return &StripeWebhookController{
		usecase: usecase,
	}
}

func (u *StripeWebhookController) Handle(fiberCtx *fiber.Ctx, req *StripeWebhookRequest) (*StripeWebhookResponse, error) {

	payload := fiberCtx.Body()

	err := u.usecase.Execute(fiberCtx.UserContext(), payload, req.Signature)
	if err != nil {
		return nil, err
	}

	return &StripeWebhookResponse{Message: "Order created successfully"}, nil

}
