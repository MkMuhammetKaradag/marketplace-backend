package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrUnauthorized       = errors.New("unauthorized access")
)
