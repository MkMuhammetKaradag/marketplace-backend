// internal/notification-service/transport/messaging/controller/user_activation.go
package controller

import (
	"context"

	"fmt"

	"marketplace/internal/notification-service/transport/messaging/usecase"
	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type UserActivationHandler struct {
	usecase usecase.UserActivationUseCase
}

func NewUserActivationHandler(usecase usecase.UserActivationUseCase) *UserActivationHandler {
	return &UserActivationHandler{
		usecase: usecase,
	}
}

func (h *UserActivationHandler) Handle(ctx context.Context, msg *pb.Message) error {

	data := msg.GetUserActivationEmailData()
	if data == nil {
		return fmt.Errorf("payload is nil or not UserActivationEmailData for message ID: %s", msg.Id)
	}

	userEmail := data.Email
	userName := data.Username
	userActivationCode := data.ActivationCode
	

	activationIDUUID, err := uuid.Parse(data.ActivationId)
	if err != nil {
		return fmt.Errorf("invalid activation id format: %w", err)
	}
	return h.usecase.Execute(ctx, activationIDUUID, userEmail, userName, userActivationCode)
}
