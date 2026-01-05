package usecase

import (
	"context"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

type GetRecommendedProductsUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Product, error)
}

type getRecommendedProductsUseCase struct {
	productRepository domain.ProductRepository
}

func NewGetRecommendedProductsUseCase(productRepository domain.ProductRepository) GetRecommendedProductsUseCase {
	return &getRecommendedProductsUseCase{
		productRepository: productRepository,
	}
}

func (c *getRecommendedProductsUseCase) Execute(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.Product, error) {
	// Repository'deki o meşhur vektör benzerliği sorgusunu çağırıyoruz
	return c.productRepository.GetRecommendedProducts(ctx, userID, limit)
}
