package usecase

import (
	"context"
	"marketplace/internal/product-service/domain"
)

type SearchProductsUseCase interface {
	Execute(ctx context.Context, req domain.SearchProductsParams) ([]*domain.Product, error)
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

func (c *searchProductsUseCase) Execute(ctx context.Context, req domain.SearchProductsParams) ([]*domain.Product, error) {
	var vector []float32
	var err error

	if req.Query != "" {
		vector, err = c.aiProvider.GetVector(req.Query)
		if err != nil {
			return nil, err
		}
	}

	// Filtreleri repository'ye paslıyoruz
	return c.productRepository.SearchProducts(ctx, vector, req)
}
