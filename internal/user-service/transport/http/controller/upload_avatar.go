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

	userIDStr, err := uuid.Parse(fbr.Get("X-User-ID"))

	if err != nil {
		return nil, err
	}
	fileHeader, err := fbr.FormFile("avatar")
	if err != nil {
		return nil, err
	}

	// 3. Dosyayı aç

	if err := h.usecase.Execute(fbr.UserContext(), userIDStr, fileHeader); err != nil {
		return nil, err
	}

	return &UploadAvatarResponse{Message: "Avatar uploaded successfully"}, nil
}
