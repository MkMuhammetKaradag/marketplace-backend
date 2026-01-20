package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/order-service/domain"

	"github.com/google/uuid"
)

type CreateOrderUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}

type createOrderUseCase struct {
	basketRepository  domain.OrderRepository
	grpcProductClient domain.ProductClient
	grpcBasketClient  domain.BasketClient
}

func NewCreateOrderUseCase(basketRepository domain.OrderRepository, grpcProductClient domain.ProductClient, grpcBasketClient domain.BasketClient) CreateOrderUseCase {
	return &createOrderUseCase{
		basketRepository:  basketRepository,
		grpcProductClient: grpcProductClient,
		grpcBasketClient:  grpcBasketClient,
	}
}

func (u *createOrderUseCase) Execute(ctx context.Context, userID uuid.UUID) error {

	basket, err := u.grpcBasketClient.GetBasket(ctx, userID.String())
	if err != nil {
		return err
	}
	if basket == nil {
		return fmt.Errorf("basket not found")
	}

	var ids []string
	for _, item := range basket.Items {
		ids = append(ids, item.ProductId)
	}
	productResponse, err := u.grpcProductClient.GetProductsByIds(ctx, ids)
	if err != nil {
		return err
	}
	if productResponse == nil {
		return fmt.Errorf("products not found")
	}
	fmt.Println("productResponse:", productResponse)
	return nil
}
