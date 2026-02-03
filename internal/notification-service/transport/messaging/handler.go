package messaginghandler

import (
	"marketplace/internal/notification-service/domain"
	"marketplace/internal/notification-service/transport/messaging/controller"
	"marketplace/internal/notification-service/transport/messaging/usecase"

	pb "marketplace/pkg/proto/events"
)

// type Handlers struct {
// 	emailProvider domain.EmailProvider
// }

// func NewMessageHandlers(email domain.EmailProvider) *Handlers {
// 	return &Handlers{
// 		emailProvider: email,
// 	}
// }

func SetupMessageHandlers(email domain.EmailProvider) map[pb.MessageType]domain.MessageHandler {

	userActivatiomUseCase := usecase.NewUserActivationUseCase(email)
	userActivationHandler := controller.NewUserActivationHandler(userActivatiomUseCase)

	return map[pb.MessageType]domain.MessageHandler{
		pb.MessageType_USER_ACTIVATION_EMAIL: userActivationHandler,
	}
}
