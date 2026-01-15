package controller

import (
	"marketplace/internal/basket-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ClearBasketRequest struct {
}

type ClearBasketResponse struct {
	Message string `json:"message"`
}
type ClearBasketController struct {
	usecase usecase.ClearBasketUseCase
}

func NewClearBasketController(usecase usecase.ClearBasketUseCase) *ClearBasketController {
	return &ClearBasketController{
		usecase: usecase,
	}
}

func (c *ClearBasketController) Handle(fiberCtx *fiber.Ctx, req *ClearBasketRequest) (*ClearBasketResponse, error) {

	userId, err := uuid.Parse(fiberCtx.Get("X-User-ID"))
	if err != nil {
		return nil, err
	}

	err = c.usecase.Execute(fiberCtx.UserContext(), userId)
	if err != nil {
		return nil, err
	}

	return &ClearBasketResponse{Message: "Basket cleared successfully"}, nil

}
