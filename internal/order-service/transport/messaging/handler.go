// internal/order-service/transport/messaging/handler.go
package messaginghandler

import (
	"marketplace/internal/order-service/domain"
	"marketplace/internal/order-service/transport/messaging/controller"
	"marketplace/internal/order-service/transport/messaging/usecase"

	eventsProto "marketplace/pkg/proto/events"
)

type Handlers struct {
	PaymentSuccess domain.MessageHandler
	PaymentFailure domain.MessageHandler
}

func NewMessageHandlers(repository domain.OrderRepository) *Handlers {
	return &Handlers{
		PaymentSuccess: controller.NewPaymentSuccessHandler(
			usecase.NewPaymentSuccessUseCase(repository),
		),
		PaymentFailure: controller.NewPaymentFailureHandler(
			usecase.NewPaymentFailureUseCase(repository),
		),
	}
}

func SetupMessageHandlers(repository domain.OrderRepository) map[eventsProto.MessageType]domain.MessageHandler {
	h := NewMessageHandlers(repository)
	return map[eventsProto.MessageType]domain.MessageHandler{
		eventsProto.MessageType_PAYMENT_SUCCESSFUL: h.PaymentSuccess,
		eventsProto.MessageType_PAYMENT_FAILED:     h.PaymentFailure,
	}
}
