package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/basket-service/domain"
)

type ProductPriceUpdatedUseCase interface {
	Execute(ctx context.Context, productId string, price float64) error
}
type productPriceUpdatedUseCase struct {
	repository domain.BasketRedisRepository
}

func NewProductPriceUpdatedUseCase(repository domain.BasketRedisRepository) ProductPriceUpdatedUseCase {
	return &productPriceUpdatedUseCase{
		repository: repository,
	}
}

func (u *productPriceUpdatedUseCase) Execute(ctx context.Context, productId string, price float64) error {

	err := u.repository.UpdateProductPriceInAllBaskets(ctx, productId, price)
	if err != nil {
		return fmt.Errorf("failed to update product price: %w", err)
	}

	return nil
}
