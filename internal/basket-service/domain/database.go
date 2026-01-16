// internal/basket-service/domain/database.go
package domain

import "context"

type BasketPostgresRepository interface {
	Close() error
}
type BasketRedisRepository interface {
	Close() error
	GetBasket(ctx context.Context, userID string) (*Basket, error)
	UpdateBasket(ctx context.Context, basket *Basket) error
	ClearBasket(ctx context.Context, userID string) error
	UpdateProductPriceInAllBaskets(ctx context.Context, productID string, newPrice float64) error
}
