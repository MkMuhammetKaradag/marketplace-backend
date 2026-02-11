package controller

import (
	"marketplace/internal/user-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadAvatarRequest struct {
}

type UploadAvatarController struct {
	usecase usecase.UploadAvatarUseCase
}

type UploadAvatarResponse struct {
	Message string `json:"message"`
}

func NewUploadAvatarController(usecase usecase.UploadAvatarUseCase) *UploadAvatarController {
	return &UploadAvatarController{
		usecase: usecase,
	}
}

func (h *UploadAvatarController) Handle(fbr *fiber.Ctx, req *UploadAvatarRequest) (*UploadAvatarResponse, error) {

	userIDStr := fbr.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid or missing user id")
	}

	fileHeader, err := fbr.FormFile("avatar")
	if err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := h.usecase.Execute(fbr.UserContext(), userID, file); err != nil {
		return nil, err
	}

	return &UploadAvatarResponse{Message: "Avatar uploaded successfully"}, nil
}
