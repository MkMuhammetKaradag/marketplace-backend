// internal/user-service/transport/http/controller/signup.go
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

// Handle godoc
// @Summary Register a new user
// @Description Creates a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body SignUpRequest true "Sign Up Request"
// @Success 200 {object} SignUpResponse
// @Router /users/signup [post]
func (h *SignUpController) Handle(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error) {
	err := h.usecase.Execute(ctx, &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &SignUpResponse{Message: " Please check your email"}, nil
}
