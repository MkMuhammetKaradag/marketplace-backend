package usecase

import (
	"context"
	"marketplace/internal/basket-service/domain"

	"github.com/google/uuid"
)

type BasketCountUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) (int, error)
}

type basketCountUseCase struct {
	basketRepository domain.BasketRedisRepository
}

func NewBasketCountUseCase(basketRepository domain.BasketRedisRepository) BasketCountUseCase {
	return &basketCountUseCase{
		basketRepository: basketRepository,
	}
}

func (u *basketCountUseCase) Execute(ctx context.Context, userID uuid.UUID) (int, error) {

	basket, err := u.basketRepository.GetBasket(ctx, userID.String())
	if err != nil {
		return 0, err
	}

	if basket == nil {
		return 0, nil
	}

	totalCount := 0
	for _, item := range basket.Items {
		totalCount += item.Quantity
	}

	return totalCount, nil
}
