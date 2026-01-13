// internal/basket-service/domain/database.go
package domain

type BasketRepository interface {
	Close() error
}
