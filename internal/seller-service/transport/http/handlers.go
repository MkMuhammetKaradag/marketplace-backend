package http

import (
	"marketplace/internal/seller-service/domain"
	"marketplace/internal/seller-service/transport/http/controller"
	"marketplace/internal/seller-service/transport/http/usecase"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Store *storeHandlers
	Admin *adminHandlers
}

type storeHandlers struct {
	Create       *controller.CreateSellerController
	GetByUserID  *controller.GetSellerByUserIDController
	UploadLogo   *controller.UploadStoreLogoController
	UploadBanner *controller.UploadStoreBannerController
}

type adminHandlers struct {
	Approve *controller.ApproveSellerController
	Reject  *controller.RejectSellerController
}

func NewHandlers(
	repository domain.SellerRepository,
	messaging domain.Messaging,
	cloudinary domain.ImageService,
) *Handlers {
	return &Handlers{
		Store: &storeHandlers{
			Create:       controller.NewCreateSellerController(usecase.NewCreateSellerUseCase(repository)),
			GetByUserID:  controller.NewGetSellerByUserIDController(usecase.NewGetSellerByUserIDUseCase(repository)),
			UploadLogo:   controller.NewUploadStoreLogoController(usecase.NewUploadStoreLogoUseCase(repository, cloudinary)),
			UploadBanner: controller.NewUploadStoreBannerController(usecase.NewUploadStoreBannerUseCase(repository, cloudinary)),
		},
		Admin: &adminHandlers{
			Approve: controller.NewApproveSellerController(usecase.NewApproveSellerUseCase(repository, messaging)),
			Reject:  controller.NewRejectSellerController(usecase.NewRejectSellerUseCase(repository, messaging)),
		},
	}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "hello from seller service",
		"info":    "Fiber handler connected to domain layer",
	})
}
