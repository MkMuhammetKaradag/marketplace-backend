// internal/user-service/transport/messaging/handler.go
package messaginghandler

import (
	"marketplace/internal/user-service/domain"
	"marketplace/internal/user-service/transport/messaging/controller"
	"marketplace/internal/user-service/transport/messaging/usecase"

	pb "marketplace/pkg/proto/events"
)

type Handlers struct {
	SellerApproved domain.MessageHandler
}

func NewHandlers(repo domain.UserRepository) *Handlers {

	return &Handlers{
		SellerApproved: controller.NewSellerApprovedHandler(
			usecase.NewSellerApprovedUseCase(repo),
		),
	}
}

func SetupMessageHandlers(repo domain.UserRepository) map[pb.MessageType]domain.MessageHandler {
	h := NewHandlers(repo)

	return map[pb.MessageType]domain.MessageHandler{
		pb.MessageType_SELLER_APPROVED: h.SellerApproved,
	}
}
