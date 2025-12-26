package controller

import (
	"marketplace/internal/seller-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadStoreLogoRequest struct {
	SellerID uuid.UUID `params:"seller_id"`
}

type UploadStoreLogoController struct {
	usecase usecase.UploadStoreLogoUseCase
}

type UploadStoreLogoResponse struct {
	Message string `json:"message"`
}

func NewUploadStoreLogoController(usecase usecase.UploadStoreLogoUseCase) *UploadStoreLogoController {
	return &UploadStoreLogoController{
		usecase: usecase,
	}
}

func (h *UploadStoreLogoController) Handle(fbr *fiber.Ctx, req *UploadStoreLogoRequest) (*UploadStoreLogoResponse, error) {

	userID, err := uuid.Parse(fbr.Get("X-User-ID"))

	if err != nil {
		return nil, err
	}
	fileHeader, err := fbr.FormFile("store_logo")
	if err != nil {
		return nil, err
	}

	if err := h.usecase.Execute(fbr.UserContext(), userID, req.SellerID, fileHeader); err != nil {
		return nil, err
	}

	return &UploadStoreLogoResponse{Message: "Store logo uploaded successfully"}, nil
}
