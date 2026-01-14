package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/basket-service/domain"

	"github.com/google/uuid"
)

type RemoveItemUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
}

type removeItemUseCase struct {
	basketRepository domain.BasketRedisRepository
}

func NewRemoveItemUseCase(basketRepository domain.BasketRedisRepository) RemoveItemUseCase {
	return &removeItemUseCase{
		basketRepository: basketRepository,
	}
}

func (u *removeItemUseCase) Execute(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	basket, err := u.basketRepository.GetBasket(ctx, userID.String())
	if err != nil {
		return err
	}

	if basket == nil || len(basket.Items) == 0 {
		return fmt.Errorf("basket not found or already empty")
	}

	newItems := []domain.BasketItem{}
	found := false
	for _, item := range basket.Items {
		if item.ProductID != productID {
			newItems = append(newItems, item)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("product not found in basket")
	}

	basket.Items = newItems

	return u.basketRepository.UpdateBasket(ctx, basket)
}
