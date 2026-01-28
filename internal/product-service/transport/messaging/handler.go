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

	// orderCreatedUseCase := usecase.NewOrderCreatedUseCase(repository)
	// orderCreatedHandler := controller.NewOrderCreatedHandler(orderCreatedUseCase)

	paymentSuccessUseCase := usecase.NewPaymentSuccessUseCase(repository)
	paymentSuccessHandler := controller.NewPaymentSuccessHandler(paymentSuccessUseCase)

	return map[pb.MessageType]domain.MessageHandler{
		pb.MessageType_SELLER_APPROVED:    sellerApprovedHandler,
		pb.MessageType_USER_CREATED:       userCreatedHandler,
		// pb.MessageType_ORDER_CREATED:      orderCreatedHandler,
		pb.MessageType_PAYMENT_SUCCESSFUL: paymentSuccessHandler,
	}
}
