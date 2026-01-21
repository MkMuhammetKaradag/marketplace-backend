package domain

import "context"

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	Close() error
}
