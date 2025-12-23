// internal/user-service/transport/http/handlers.go
package http

import (
	"github.com/gofiber/fiber/v2"

	"marketplace/internal/seller-service/domain"
	"marketplace/internal/seller-service/transport/http/controller"
	"marketplace/internal/seller-service/transport/http/usecase"
)

type Handlers struct {
	sellerRepository domain.SellerRepository
	kafkaMessaging   domain.Messaging
}

func NewHandlers(repository domain.SellerRepository, messaging domain.Messaging) *Handlers {
	return &Handlers{sellerRepository: repository, kafkaMessaging: messaging}
}

func (h *Handlers) Hello(c *fiber.Ctx) error {

	resp := HelloResponse{
		Message: "hhelu seller service",
		Info:    "Fiber handler connected to domain layer",
	}
	return c.JSON(resp)
}

func (h *Handlers) CreateSeller() *controller.CreateSellerController {
	createSellerUseCase := usecase.NewCreateSellerUseCase(h.sellerRepository)
	return controller.NewCreateSellerController(createSellerUseCase)
}

func (h *Handlers) RejectSeller() *controller.RejectSellerController {
	rejectSellerUseCase := usecase.NewRejectSellerUseCase(h.sellerRepository, h.kafkaMessaging)
	return controller.NewRejectSellerController(rejectSellerUseCase)
}

func (h *Handlers) ApproveSeller() *controller.ApproveSellerController {
	approveSellerUseCase := usecase.NewApproveSellerUseCase(h.sellerRepository, h.kafkaMessaging)
	return controller.NewApproveSellerController(approveSellerUseCase)
}

func (h *Handlers) GetSellerByUserID() *controller.GetSellerByUserIDController {
	getSellerByUserIDUseCase := usecase.NewGetSellerByUserIDUseCase(h.sellerRepository)
	return controller.NewGetSellerByUserIDController(getSellerByUserIDUseCase)
}

type HelloResponse struct {
	Message string `json:"message"`
	Info    string `json:"info"`
}
