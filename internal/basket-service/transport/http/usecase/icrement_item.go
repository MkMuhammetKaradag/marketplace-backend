package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/basket-service/domain"

	"github.com/google/uuid"
)

type IncrementItemUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
}

type incrementItemUseCase struct {
	basketRepository  domain.BasketRedisRepository
	grpcProductClient domain.ProductClient
}

func NewIncrementItemUseCase(basketRepository domain.BasketRedisRepository, grpcProductClient domain.ProductClient) IncrementItemUseCase {
	return &incrementItemUseCase{
		basketRepository:  basketRepository,
		grpcProductClient: grpcProductClient,
	}
}

func (u *incrementItemUseCase) Execute(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {

	product, err := u.grpcProductClient.GetProductForBasket(ctx, productID.String())
	if err != nil {
		return err
	}
	if product == nil { // || !product.IsActive
		return fmt.Errorf("product is not available")
	}

	// 2. get basket
	basket, err := u.basketRepository.GetBasket(ctx, userID.String())
	if err != nil {
		return err
	}
	if basket == nil {
		return fmt.Errorf("basket not found")
	}

	found := false
	// We are performing operations directly on basket.Items (using a pointer or index).
	for i := range basket.Items {
		if basket.Items[i].ProductID == productID {
			found = true

			// check stock
			if product.Stock < int32(basket.Items[i].Quantity+1) {
				return fmt.Errorf("insufficient stock. current stock: %d", product.Stock)
			}

			// check max quantity
			if basket.Items[i].Quantity >= 20 {
				return fmt.Errorf("You can add a maximum of 20 units of one product.")
			}

			// update
			basket.Items[i].Quantity++
			basket.Items[i].Price = product.Price // update price
			break
		}
	}

	if !found {
		return fmt.Errorf("product not found in basket")
	}

	// 3. update basket
	return u.basketRepository.UpdateBasket(ctx, basket)
}
