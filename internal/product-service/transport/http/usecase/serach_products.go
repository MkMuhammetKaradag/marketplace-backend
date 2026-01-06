package usecase

import (
	"context"
	"marketplace/internal/product-service/domain"
)

type SearchProductsUseCase interface {
	Execute(ctx context.Context, limit int, query string) ([]*domain.Product, error)
}

type searchProductsUseCase struct {
	productRepository domain.ProductRepository
	aiProvider        domain.AiProvider
}

func NewSearchProductsUseCase(productRepository domain.ProductRepository, aiProvider domain.AiProvider) SearchProductsUseCase {
	return &searchProductsUseCase{
		productRepository: productRepository,
		aiProvider:        aiProvider,
	}
}

func (c *searchProductsUseCase) Execute(ctx context.Context, limit int, query string) ([]*domain.Product, error) {

	vector, err := c.aiProvider.GetVector(query)
	if err != nil {
		return nil, err
	}
	return c.productRepository.SearchProducts(ctx, vector, query, limit)
}
