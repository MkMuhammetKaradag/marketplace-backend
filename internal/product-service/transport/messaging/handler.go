package messaginghandler

import (
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/messaging/controller"
	"marketplace/internal/product-service/transport/messaging/usecase"

	pb "marketplace/pkg/proto/events"
)

type Handlers struct {
	productRepository domain.ProductRepository
}

func NewMessageHandlers(repository domain.ProductRepository) *Handlers {
	return &Handlers{productRepository: repository}
}

func SetupMessageHandlers(repository domain.ProductRepository) map[pb.MessageType]domain.MessageHandler {
	sellerApprovedUseCase := usecase.NewSellerApprovedUseCase(repository)
	sellerApprovedHandler := controller.NewSellerApprovedHandler(sellerApprovedUseCase)

	userCreatedUseCase := usecase.NewUserCreatedUseCase(repository)
	userCreatedHandler := controller.NewUserCreatedHandler(userCreatedUseCase)
	return map[pb.MessageType]domain.MessageHandler{
		pb.MessageType_SELLER_APPROVED: sellerApprovedHandler,
		pb.MessageType_USER_CREATED:    userCreatedHandler,
	}
}
