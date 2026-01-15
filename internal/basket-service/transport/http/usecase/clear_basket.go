package usecase

import (
	"context"
	"marketplace/internal/basket-service/domain"

	"github.com/google/uuid"
)

type ClearBasketUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}

type clearBasketUseCase struct {
	basketRepository domain.BasketRedisRepository
}

func NewClearBasketUseCase(basketRepository domain.BasketRedisRepository) ClearBasketUseCase {
	return &clearBasketUseCase{
		basketRepository: basketRepository,
	}
}

func (u *clearBasketUseCase) Execute(ctx context.Context, userID uuid.UUID) error {
	return u.basketRepository.ClearBasket(ctx, userID.String())
}
