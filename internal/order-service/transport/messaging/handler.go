package messaginghandler

import (
	"marketplace/internal/order-service/domain"
	"marketplace/internal/order-service/transport/messaging/controller"
	"marketplace/internal/order-service/transport/messaging/usecase"

	eventsProto "marketplace/pkg/proto/events"
)

type Handlers struct {
	orderRepository domain.OrderRepository
}

func NewMessageHandlers(repository domain.OrderRepository) *Handlers {
	return &Handlers{orderRepository: repository}
}

func SetupMessageHandlers(repository domain.OrderRepository) map[eventsProto.MessageType]domain.MessageHandler {
	paymentSuccessUseCase := usecase.NewPaymentSuccessUseCase(repository)
	paymentSuccessHandler := controller.NewPaymentSuccessHandler(paymentSuccessUseCase)

	return map[eventsProto.MessageType]domain.MessageHandler{
		eventsProto.MessageType_PAYMENT_SUCCESSFUL: paymentSuccessHandler,
	}
}
