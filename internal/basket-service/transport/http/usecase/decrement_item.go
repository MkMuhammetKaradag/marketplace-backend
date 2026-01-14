package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/basket-service/domain"

	"github.com/google/uuid"
)

type DecrementItemUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
}

type decrementItemUseCase struct {
	basketRepository domain.BasketRedisRepository
}

func NewDecrementItemUseCase(basketRepository domain.BasketRedisRepository) DecrementItemUseCase {
	return &decrementItemUseCase{
		basketRepository: basketRepository,
	}
}

func (u *decrementItemUseCase) Execute(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	basket, err := u.basketRepository.GetBasket(ctx, userID.String())
	if err != nil {
		return err
	}
	if basket == nil {
		return fmt.Errorf("basket not found")
	}


	found := false
	newItems := []domain.BasketItem{}

	for _, item := range basket.Items {
		if item.ProductID == productID {
			found = true
			if item.Quantity > 1 {
				item.Quantity--
				newItems = append(newItems, item)
			}
		} else {
			newItems = append(newItems, item)
		}
	}

	if !found {
		return fmt.Errorf("product not found in basket")
	}

	
	basket.Items = newItems
	return u.basketRepository.UpdateBasket(ctx, basket)
}
