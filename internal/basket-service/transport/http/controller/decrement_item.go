package controller

import (
	"marketplace/internal/basket-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type DecrementItemRequest struct {
	ProductID uuid.UUID `params:"product_id"`
}
type DecrementItemResponse struct {
	Message string `json:"message"`
}
type DecrementItemController struct {
	usecase usecase.DecrementItemUseCase
}

func NewDecrementItemController(usecase usecase.DecrementItemUseCase) *DecrementItemController {
	return &DecrementItemController{
		usecase: usecase,
	}
}

func (c *DecrementItemController) Handle(fiberCtx *fiber.Ctx, req *DecrementItemRequest) (*DecrementItemResponse, error) {

	userId, err := uuid.Parse(fiberCtx.Get("X-User-ID"))
	if err != nil {
		return nil, err
	}

	err = c.usecase.Execute(fiberCtx.UserContext(), userId, req.ProductID)
	if err != nil {
		return nil, err
	}

	return &DecrementItemResponse{Message: "Item decremented from basket successfully"}, nil

}
