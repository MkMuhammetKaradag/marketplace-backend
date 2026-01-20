package controller

import (
	"marketplace/internal/order-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
}

type CreateOrderResponse struct {
	Message string `json:"message"`
}
type CreateOrderController struct {
	usecase usecase.CreateOrderUseCase
}

func NewCreateOrderController(usecase usecase.CreateOrderUseCase) *CreateOrderController {
	return &CreateOrderController{
		usecase: usecase,
	}
}

func (c *CreateOrderController) Handle(fiberCtx *fiber.Ctx, req *CreateOrderRequest) (*CreateOrderResponse, error) {

	userId, err := uuid.Parse(fiberCtx.Get("X-User-ID"))
	if err != nil {
		return nil, err
	}

	err = c.usecase.Execute(fiberCtx.UserContext(), userId)
	if err != nil {
		return nil, err
	}

	return &CreateOrderResponse{Message: "Order created successfully"}, nil

}
