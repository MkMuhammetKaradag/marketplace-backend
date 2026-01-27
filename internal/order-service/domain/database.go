package domain

import (
	"context"

	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status OrderStatus) error
	Close() error
}
