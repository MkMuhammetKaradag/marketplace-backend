package controller

import (
	"errors"
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type SearchProductsRequest struct {
	Limit      int      `json:"limit"`
	Query      string   `json:"query"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	CategoryID *string  `json:"category_id"`
}

type SearchProductsResponse struct {
	Products []*domain.Product `json:"products"`
}

type SearchProductsController struct {
	usecase usecase.SearchProductsUseCase
}

func NewSearchProductsController(usecase usecase.SearchProductsUseCase) *SearchProductsController {
	return &SearchProductsController{
		usecase: usecase,
	}
}

// Handle godoc
// @Summary Search products
// @Description Search for products using vector search
// @Tags products
// @Accept json
// @Produce json
// @Param limit query int false "Limit results"
// @Param query query string false "Search query"
// @Success 200 {object} SearchProductsResponse
// @Router /products/search [get]
func (c *SearchProductsController) Handle(fiberCtx *fiber.Ctx, req *SearchProductsRequest) (*SearchProductsResponse, error) {
	params := domain.SearchProductsParams{
		Limit:      req.Limit,
		Query:      req.Query,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		CategoryID: req.CategoryID,
	}

	products, err := c.usecase.Execute(fiberCtx.UserContext(), params)
	if err != nil {
		return nil, errors.New("failed to get search products")
	}

	return &SearchProductsResponse{
		Products: products,
	}, nil
}
