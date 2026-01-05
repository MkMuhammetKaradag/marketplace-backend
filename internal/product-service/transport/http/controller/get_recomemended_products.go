package controller

import (
	"errors"
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetRecommendationsRequest struct {
	Limit int `json:"limit"`
}

type GetRecommendationsResponse struct {
	Items []*domain.Product `json:"items"`
}

type GetRecommendationsController struct {
	usecase usecase.GetRecommendedProductsUseCase
}

func NewGetRecommendationsController(usecase usecase.GetRecommendedProductsUseCase) *GetRecommendationsController {
	return &GetRecommendationsController{
		usecase: usecase,
	}
}

func (c *GetRecommendationsController) Handle(fiberCtx *fiber.Ctx, req *GetRecommendationsRequest) (*GetRecommendationsResponse, error) {
	// 1. Header'dan UserID al
	userIDStr := fiberCtx.Get("X-User-ID")
	if userIDStr == "" {
		return nil, errors.New("user id is required")
	}

	parsedUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	// 2. UseCase'i çalıştır (Örn: Ana sayfa için 20 ürün getir)
	products, err := c.usecase.Execute(fiberCtx.UserContext(), parsedUserID, req.Limit)
	if err != nil {
		return nil, errors.New("failed to get recommended products")
	}

	// 3. Sonucu dön
	return &GetRecommendationsResponse{
		Items: products,
	}, nil
}
