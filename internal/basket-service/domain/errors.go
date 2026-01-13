// internal/basket-service/domain/errors.go
package domain

import "errors"

var (
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrBasketNotFound = errors.New("basket not found")
)
