package usecase

import (
	"context"
	"marketplace/internal/basket-service/domain"

	"github.com/google/uuid"
)

type AddItemUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, p *domain.BasketItem) error
}

type addItemUseCase struct {
	basketRepository domain.BasketRedisRepository
}

func NewAddItemUseCase(basketRepository domain.BasketRedisRepository) AddItemUseCase {
	return &addItemUseCase{
		basketRepository: basketRepository,
	}
}

func (u *addItemUseCase) Execute(ctx context.Context, userID uuid.UUID, p *domain.BasketItem) error {

	basket, err := u.basketRepository.GetBasket(ctx, userID.String())
	if err != nil {
		return err
	}

	if basket == nil {
		basket = &domain.Basket{UserID: userID, Items: []domain.BasketItem{}}
	}

	found := false
	for i, item := range basket.Items {
		if item.ProductID == p.ProductID {
			basket.Items[i].Quantity += p.Quantity
			found = true
			break
		}
	}

	if !found {
		basket.Items = append(basket.Items, *p)
	}

	return u.basketRepository.UpdateBasket(ctx, basket)
}
