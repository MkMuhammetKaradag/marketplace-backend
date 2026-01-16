package controller

import (
	"marketplace/internal/basket-service/domain"
	"marketplace/internal/basket-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetBasketRequest struct {
}

type GetBasketResponse struct {
	Basket *domain.BasketResponse `json:"basket"`
}
type GetBasketController struct {
	usecase usecase.GetBasketUseCase
}

func NewGetBasketController(usecase usecase.GetBasketUseCase) *GetBasketController {
	return &GetBasketController{
		usecase: usecase,
	}
}

func (c *GetBasketController) Handle(fiberCtx *fiber.Ctx, req *GetBasketRequest) (*GetBasketResponse, error) {

	userId, err := uuid.Parse(fiberCtx.Get("X-User-ID"))
	if err != nil {
		return nil, err
	}

	basket, err := c.usecase.Execute(fiberCtx.UserContext(), userId)
	if err != nil {
		return nil, err
	}

	return &GetBasketResponse{Basket: basket}, nil

}
