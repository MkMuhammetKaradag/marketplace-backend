package usecase

import (
	"context"
	"marketplace/internal/product-service/domain"
)

type CreateCategoryUseCase interface {
	Execute(ctx context.Context, req *domain.Category) error
}

type createCategoryUseCase struct {
	categoryRepository domain.ProductRepository
}

func NewCreateCategoryUseCase(categoryRepository domain.ProductRepository) CreateCategoryUseCase {
	return &createCategoryUseCase{
		categoryRepository: categoryRepository,
	}
}

func (c *createCategoryUseCase) Execute(ctx context.Context, req *domain.Category) error {

	return c.categoryRepository.CreateCategory(ctx, req)
}
