package controller

import (
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetFavoritesRequest struct {
}

type GetFavoritesResponse struct {
	Products []*domain.FavoriteItem `json:"products"`
}

type GetFavoritesController struct {
	usecase usecase.GetFavoritesUseCase
}

func NewGetFavoritesController(usecase usecase.GetFavoritesUseCase) *GetFavoritesController {
	return &GetFavoritesController{
		usecase: usecase,
	}
}

// Handle godoc
// @Summary Get favorites
// @Description Get  user favorites
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} GetFavoritesResponse
// @Router /products/favorites [get]
func (c *GetFavoritesController) Handle(fiberCtx *fiber.Ctx, req *GetFavoritesRequest) (*GetFavoritesResponse, error) {

	userIDStr := fiberCtx.Get("X-User-ID")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	products, err := c.usecase.Execute(fiberCtx.UserContext(), userID)
	if err != nil {
		return nil, err
	}

	return &GetFavoritesResponse{
		Products: products,
	}, nil
}
