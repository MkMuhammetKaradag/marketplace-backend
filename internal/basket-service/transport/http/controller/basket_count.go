package controller

import (
	"marketplace/internal/basket-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BasketCountRequest struct {
}

type BasketCountResponse struct {
	Count int `json:"count"`
}
type BasketCountController struct {
	usecase usecase.BasketCountUseCase
}

func NewBasketCountController(usecase usecase.BasketCountUseCase) *BasketCountController {
	return &BasketCountController{
		usecase: usecase,
	}
}

func (c *BasketCountController) Handle(fiberCtx *fiber.Ctx, req *BasketCountRequest) (*BasketCountResponse, error) {

	userId, err := uuid.Parse(fiberCtx.Get("X-User-ID"))
	if err != nil {
		return nil, err
	}

	count, err := c.usecase.Execute(fiberCtx.UserContext(), userId)
	if err != nil {
		return nil, err
	}

	return &BasketCountResponse{Count: count}, nil

}
