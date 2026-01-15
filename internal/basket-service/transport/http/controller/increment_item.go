package controller

import (
	"marketplace/internal/basket-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IncrementItemRequest struct {
	ProductID uuid.UUID `params:"product_id"`
}
type IncrementItemResponse struct {
	Message string `json:"message"`
}
type IncrementItemController struct {
	usecase usecase.IncrementItemUseCase
}

func NewIncrementItemController(usecase usecase.IncrementItemUseCase) *IncrementItemController {
	return &IncrementItemController{
		usecase: usecase,
	}
}

func (c *IncrementItemController) Handle(fiberCtx *fiber.Ctx, req *IncrementItemRequest) (*IncrementItemResponse, error) {

	userId, err := uuid.Parse(fiberCtx.Get("X-User-ID"))
	if err != nil {
		return nil, err
	}

	err = c.usecase.Execute(fiberCtx.UserContext(), userId, req.ProductID)
	if err != nil {
		return nil, err
	}

	return &IncrementItemResponse{Message: "Item incremented from basket successfully"}, nil

}
