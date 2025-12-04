// internal/user-service/domain/errors.go
package domain

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized access")
)
