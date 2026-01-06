package controller

import (
	"errors"
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type SearchProductsRequest struct {
	Limit int    `json:"limit"`
	Query string `json:"query"`
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

func (c *SearchProductsController) Handle(fiberCtx *fiber.Ctx, req *SearchProductsRequest) (*SearchProductsResponse, error) {

	products, err := c.usecase.Execute(fiberCtx.UserContext(), req.Limit, req.Query)
	if err != nil {
		return nil, errors.New("failed to get search products")
	}

	return &SearchProductsResponse{
		Products: products,
	}, nil
}
