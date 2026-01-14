package controller

import (
	"marketplace/internal/basket-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RemoveItemRequest struct {
	ProductID uuid.UUID `params:"product_id"`
}

type RemoveItemResponse struct {
	Message string `json:"message"`
}
type RemoveItemController struct {
	usecase usecase.RemoveItemUseCase
}

func NewRemoveItemController(usecase usecase.RemoveItemUseCase) *RemoveItemController {
	return &RemoveItemController{
		usecase: usecase,
	}
}

func (c *RemoveItemController) Handle(fiberCtx *fiber.Ctx, req *RemoveItemRequest) (*RemoveItemResponse, error) {

	userId, err := uuid.Parse(fiberCtx.Get("X-User-ID"))
	if err != nil {
		return nil, err
	}

	err = c.usecase.Execute(fiberCtx.UserContext(), userId, req.ProductID)
	if err != nil {
		return nil, err
	}

	return &RemoveItemResponse{Message: "Item removed from basket successfully"}, nil

}
