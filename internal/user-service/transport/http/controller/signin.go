// internal/user-service/transport/http/controller/signIn.go
package controller

import (
	"marketplace/internal/user-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type SignInRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=8"`
}

type SignInResponse struct {
	Message string `json:"message"`
}
type SignInController struct {
	usecase usecase.SignInUseCase
}

func NewSignInController(usecase usecase.SignInUseCase) *SignInController {
	return &SignInController{
		usecase: usecase,
	}
}

func (h *SignInController) Handle(fiberCtx *fiber.Ctx, req *SignInRequest) (*SignInResponse, error) {
	err := h.usecase.Execute(fiberCtx, req.Identifier, req.Password)
	if err != nil {
		return nil, err
	}

	return &SignInResponse{Message: "Signin Success"}, nil
}
