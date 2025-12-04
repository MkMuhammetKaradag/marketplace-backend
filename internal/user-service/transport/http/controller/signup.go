package controller

import (
	"context"
	"marketplace/internal/user-service/domain"
	"marketplace/internal/user-service/transport/http/usecase"
)

type SignUpRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignUpResponse struct {
	Message string `json:"message"`
}
type SignUpController struct {
	usecase usecase.SignUpUseCase
}

func NewSignUpController(usecase usecase.SignUpUseCase) *SignUpController {
	return &SignUpController{
		usecase: usecase,
	}
}

func (h *SignUpController) Handle(ctx context.Context, req *SignUpRequest) (*SignUpResponse, int, error) {
	status, err := h.usecase.Execute(ctx, &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, status, err
	}

	return &SignUpResponse{Message: " Please check your email"}, status, nil
}
