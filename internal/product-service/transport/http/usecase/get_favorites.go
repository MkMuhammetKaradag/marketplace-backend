package usecase

import (
	"context"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

type GetFavoritesUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) ([]*domain.FavoriteItem, error)
}

type getFavoritesUseCase struct {
	productRepository domain.ProductRepository
}

func NewGetFavoritesUseCase(productRepository domain.ProductRepository) GetFavoritesUseCase {
	return &getFavoritesUseCase{
		productRepository: productRepository,
	}
}

func (c *getFavoritesUseCase) Execute(ctx context.Context, userID uuid.UUID) ([]*domain.FavoriteItem, error) {

	products, err := c.productRepository.GetUserFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}

	return products, nil
}
