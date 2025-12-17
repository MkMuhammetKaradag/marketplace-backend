package messaginghandler

import (
	"marketplace/internal/user-service/domain"
	"marketplace/internal/user-service/transport/messaging/controller"
	"marketplace/internal/user-service/transport/messaging/usecase"

	pb "marketplace/pkg/proto/events"
)

type Handlers struct {
	userRepository domain.UserRepository
}

func NewMessageHandlers(repository domain.UserRepository) *Handlers {
	return &Handlers{userRepository: repository}
}

func SetupMessageHandlers(repository domain.UserRepository) map[pb.MessageType]domain.MessageHandler {
	sellerApprovedUseCase := usecase.NewSellerApprovedUseCase(repository)
	sellerApprovedHandler := controller.NewSellerApprovedHandler(sellerApprovedUseCase)

	return map[pb.MessageType]domain.MessageHandler{
		pb.MessageType_SELLER_APPROVED: sellerApprovedHandler,
	}
}
