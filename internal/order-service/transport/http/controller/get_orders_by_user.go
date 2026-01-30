package controller

import (
	"fmt"
	"marketplace/internal/order-service/domain"
	"marketplace/internal/order-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetOrdersByUserRequest struct {
}

type GetOrdersByUserResponse struct {
	Message string         `json:"message"`
	Orders  []domain.Order `json:"orders"`
}
type GetOrdersByUserController struct {
	usecase usecase.GetOrderByUserUseCase
}

func NewGetOrdersByUserController(usecase usecase.GetOrderByUserUseCase) *GetOrdersByUserController {
	return &GetOrdersByUserController{
		usecase: usecase,
	}
}

func (c *GetOrdersByUserController) Handle(fiberCtx *fiber.Ctx, req *GetOrdersByUserRequest) (*GetOrdersByUserResponse, error) {

	userId, err := uuid.Parse(fiberCtx.Get("X-User-ID"))
	if err != nil {
		fmt.Println("get orders by user error", err)
		return nil, err
	}

	orders, err := c.usecase.Execute(fiberCtx.UserContext(), userId)
	if err != nil {
		return nil, err
	}

	return &GetOrdersByUserResponse{Message: "Order created successfully", Orders: orders}, nil

}
