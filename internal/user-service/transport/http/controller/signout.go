// internal/user-service/transport/http/controller/signout.go
package controller

import (
	"marketplace/internal/user-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type SignOutRequest struct{}
type SignOutResponse struct {
	Message string `json:"message"`
}
type SignOutController struct {
	usecase usecase.SignOutUseCase
}

func NewSignOutController(usecase usecase.SignOutUseCase) *SignOutController {
	return &SignOutController{
		usecase: usecase,
	}
}

func (h *SignOutController) Handle(fbrCtx *fiber.Ctx, req *SignOutRequest) (*SignOutResponse, error) {
	err := h.usecase.Execute(fbrCtx)
	if err != nil {
		return nil, err
	}

	return &SignOutResponse{Message: "logout successfully"}, nil
}
