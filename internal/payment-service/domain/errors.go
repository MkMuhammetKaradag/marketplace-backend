// internal/payment-service/domain/errors.go
package domain

import "errors"

var (
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrOrderNotFound = errors.New("order not found")
)
