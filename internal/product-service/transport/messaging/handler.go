// internal/product-service/transport/messaging/handler.go
package messaginghandler

import (
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/transport/messaging/controller"
	"marketplace/internal/product-service/transport/messaging/usecase"

	pb "marketplace/pkg/proto/events"
)

type Handlers struct {
	SellerApproved domain.MessageHandler
	UserCreated    domain.MessageHandler
	PaymentSuccess domain.MessageHandler
	PaymentFailure domain.MessageHandler
}

func NewMessageHandlers(repo domain.ProductRepository) *Handlers {
	return &Handlers{
		SellerApproved: controller.NewSellerApprovedHandler(
			usecase.NewSellerApprovedUseCase(repo),
		),
		UserCreated: controller.NewUserCreatedHandler(
			usecase.NewUserCreatedUseCase(repo),
		),
		PaymentSuccess: controller.NewPaymentSuccessHandler(
			usecase.NewPaymentSuccessUseCase(repo),
		),
		PaymentFailure: controller.NewPaymentFailureHandler(
			usecase.NewPaymentFailureUseCase(repo),
		),
	}
}

func SetupMessageHandlers(repo domain.ProductRepository) map[pb.MessageType]domain.MessageHandler {
	h := NewMessageHandlers(repo)

	return map[pb.MessageType]domain.MessageHandler{
		pb.MessageType_SELLER_APPROVED:    h.SellerApproved,
		pb.MessageType_USER_CREATED:       h.UserCreated,
		pb.MessageType_PAYMENT_SUCCESSFUL: h.PaymentSuccess,
		pb.MessageType_PAYMENT_FAILED:     h.PaymentFailure,
	}
}
