package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/order-service/domain"

	"github.com/google/uuid"
)

type GetOrderByUserUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) ([]domain.Order, error)
}

type getOrderByUserUseCase struct {
	orderRepository domain.OrderRepository
}

func NewGetOrderByUserUseCase(orderRepository domain.OrderRepository) GetOrderByUserUseCase {
	return &getOrderByUserUseCase{
		orderRepository: orderRepository,
	}

}

func (u *getOrderByUserUseCase) Execute(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {

	orders, err := u.orderRepository.GetOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}
	return orders, nil
}
