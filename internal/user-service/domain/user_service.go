// internal/user-service/domain/user_service.go
package domain

import (
	"context"
)

type UserService interface {
	Greeting(ctx context.Context) string
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Greeting(_ context.Context) string {
	return "Hello from User Service!"
}
