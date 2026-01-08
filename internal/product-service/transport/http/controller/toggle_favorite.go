package controller

import (
	"marketplace/internal/product-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ToggleFavoriteRequest struct {
	ProductID uuid.UUID `params:"product_id"`
}

type ToggleFavoriteResponse struct {
	Message string `json:"message"`
}
type ToggleFavoriteController struct {
	usecase usecase.ToggleFavoriteUseCase
}

func NewToggleFavoriteController(usecase usecase.ToggleFavoriteUseCase) *ToggleFavoriteController {
	return &ToggleFavoriteController{
		usecase: usecase,
	}
}

func (c *ToggleFavoriteController) Handle(fiberCtx *fiber.Ctx, req *ToggleFavoriteRequest) (*ToggleFavoriteResponse, error) {

	userId := fiberCtx.Get("X-USER-ID")

	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	err = c.usecase.Execute(fiberCtx.UserContext(), userIdUUID, req.ProductID)
	if err != nil {
		return nil, err
	}

	return &ToggleFavoriteResponse{Message: "Favorite toggled successfully"}, nil

}
