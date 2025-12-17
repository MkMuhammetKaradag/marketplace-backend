package messaginghandler

import (
	"marketplace/internal/user-service/domain"
	"marketplace/internal/user-service/transport/messaging/controller"
	"marketplace/internal/user-service/transport/messaging/usecase"
	"marketplace/pkg/messaging"
)

type Handlers struct {
	userRepository domain.UserRepository
}

func NewMessageHandlers(repository domain.UserRepository) *Handlers {
	return &Handlers{userRepository: repository}
}

func SetupMessageHandlers(repository domain.UserRepository) map[messaging.MessageType]domain.MessageHandler {
	sellerApprovedUseCase := usecase.NewSellerApprovedUseCase(repository)
	sellerApprovedHandler := controller.NewSellerApprovedHandler(sellerApprovedUseCase)

	return map[messaging.MessageType]domain.MessageHandler{
		messaging.MessageType_SELLER_APPROVED: sellerApprovedHandler,
	}
}
