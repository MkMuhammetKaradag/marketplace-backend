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

func SetupMessageHandlers(email domain.EmailProvider, repository domain.NotificationRepository) map[pb.MessageType]domain.MessageHandler {

	userActivatiomUseCase := usecase.NewUserActivationUseCase(email)
	userActivationHandler := controller.NewUserActivationHandler(userActivatiomUseCase)

	userCreatedUseCase := usecase.NewUserCreatedUseCase(repository)
	userCreatedHandler := controller.NewUserCreatedHandler(userCreatedUseCase)

	orderCreatedUseCase := usecase.NewOrderCreatedUseCase(email, repository)
	orderCreatedHandler := controller.NewOrderCreatedHandler(orderCreatedUseCase)

	paymentSuccessUseCase := usecase.NewPaymentSuccessUseCase(repository, email)
	paymentSuccessHandler := controller.NewPaymentSuccessHandler(paymentSuccessUseCase)

	paymentFailedUseCase := usecase.NewPaymentFailedUseCase(repository, email)
	paymentFailedHandler := controller.NewPaymentFailedHandler(paymentFailedUseCase)

	rejectSellerUseCase := usecase.NewRejectSellerUseCase(email, repository)
	rejectSellerHandler := controller.NewRejectSellerHandler(rejectSellerUseCase)

	approveSellerUseCase := usecase.NewApproveSellerUseCase(email, repository)
	approveSellerHandler := controller.NewApproveSellerHandler(approveSellerUseCase)

	forgotPasswordUseCase := usecase.NewForgotPasswordUseCase(email, repository)
	forgotPasswordHandler := controller.NewForgotPasswordHandler(forgotPasswordUseCase)
	return map[pb.MessageType]domain.MessageHandler{
		pb.MessageType_USER_ACTIVATION_EMAIL: userActivationHandler,
		pb.MessageType_USER_CREATED:          userCreatedHandler,
		pb.MessageType_ORDER_CREATED:         orderCreatedHandler,
		pb.MessageType_PAYMENT_SUCCESSFUL:    paymentSuccessHandler,
		pb.MessageType_PAYMENT_FAILED:        paymentFailedHandler,
		pb.MessageType_SELLER_REJECTED:       rejectSellerHandler,
		pb.MessageType_USER_FORGOT_PASSWORD:  forgotPasswordHandler,
		pb.MessageType_SELLER_APPROVED:       approveSellerHandler,
	}
}
