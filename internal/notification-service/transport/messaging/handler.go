package messaginghandler

import (
	"marketplace/internal/notification-service/domain"
	"marketplace/internal/notification-service/transport/messaging/controller"
	"marketplace/internal/notification-service/transport/messaging/usecase"
	pb "marketplace/pkg/proto/events"
)

type registry struct {
	email       domain.EmailProvider
	templateMgr domain.TemplateManager
	repo        domain.NotificationRepository
	handlers    map[pb.MessageType]domain.MessageHandler
}

func SetupMessageHandlers(email domain.EmailProvider, templateMgr domain.TemplateManager, repository domain.NotificationRepository) map[pb.MessageType]domain.MessageHandler {
	r := &registry{
		email:       email,
		repo:        repository,
		templateMgr: templateMgr,
		handlers:    make(map[pb.MessageType]domain.MessageHandler),
	}

	//Grouped registration processes
	r.registerUserHandlers()
	r.registerOrderHandlers()
	r.registerPaymentHandlers()
	r.registerSellerHandlers()

	return r.handlers
}

func (r *registry) registerUserHandlers() {
	// User Activation
	activationUC := usecase.NewUserActivationUseCase(r.email, r.templateMgr)
	r.handlers[pb.MessageType_USER_ACTIVATION_EMAIL] = controller.NewUserActivationHandler(activationUC)

	// User Created (Sync)
	createdUC := usecase.NewUserCreatedUseCase(r.repo)
	r.handlers[pb.MessageType_USER_CREATED] = controller.NewUserCreatedHandler(createdUC)

	// Forgot Password
	forgotUC := usecase.NewForgotPasswordUseCase(r.email, r.repo, r.templateMgr)
	r.handlers[pb.MessageType_USER_FORGOT_PASSWORD] = controller.NewForgotPasswordHandler(forgotUC)
}

func (r *registry) registerOrderHandlers() {
	orderUC := usecase.NewOrderCreatedUseCase(r.email, r.repo, r.templateMgr)
	r.handlers[pb.MessageType_ORDER_CREATED] = controller.NewOrderCreatedHandler(orderUC)
}

func (r *registry) registerPaymentHandlers() {
	// Success
	successUC := usecase.NewPaymentSuccessUseCase(r.repo, r.email, r.templateMgr)
	r.handlers[pb.MessageType_PAYMENT_SUCCESSFUL] = controller.NewPaymentSuccessHandler(successUC)

	// Failed
	failedUC := usecase.NewPaymentFailedUseCase(r.repo, r.email, r.templateMgr)
	r.handlers[pb.MessageType_PAYMENT_FAILED] = controller.NewPaymentFailedHandler(failedUC)
}

func (r *registry) registerSellerHandlers() {
	// Reject
	rejectUC := usecase.NewRejectSellerUseCase(r.email, r.repo, r.templateMgr)
	r.handlers[pb.MessageType_SELLER_REJECTED] = controller.NewRejectSellerHandler(rejectUC)

	// Approve
	approveUC := usecase.NewApproveSellerUseCase(r.email, r.templateMgr, r.repo)
	r.handlers[pb.MessageType_SELLER_APPROVED] = controller.NewApproveSellerHandler(approveUC)
}
