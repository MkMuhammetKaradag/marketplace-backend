package controller

import (
	"marketplace/internal/user-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CreateRoleRequest struct {
	Permissions int64  `json:"permissions"`
	Name        string `json:"name"`
}
type CreateRoleController struct {
	usecase usecase.CreateRoleUsecase
}

type CreateRoleResponse struct {
	Message string `json:"message"`
}

func NewCreateRoleController(usecase usecase.CreateRoleUsecase) *CreateRoleController {
	return &CreateRoleController{
		usecase: usecase,
	}
}

func (h *CreateRoleController) Handle(fiberCtx *fiber.Ctx, req *CreateRoleRequest) (*CreateRoleResponse, error) {
	userIDSrt := fiberCtx.Get("X-User-ID")
	parsedID, err := uuid.Parse(userIDSrt)
	if err != nil {
		parsedID = uuid.Nil
	}
	err = h.usecase.Execute(fiberCtx, parsedID, req.Permissions, req.Name)
	if err != nil {
		return nil, err
	}

	return &CreateRoleResponse{Message: "Role created successfully"}, nil
}
