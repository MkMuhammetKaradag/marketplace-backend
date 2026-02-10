package controller

import (
	"context"
	"fmt"
	"marketplace/internal/notification-service/transport/messaging/usecase"
	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type ForgotPasswordHandler struct {
	usecase usecase.ForgotPasseordUseCase
}

func NewForgotPasswordHandler(usecase usecase.ForgotPasseordUseCase) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		usecase: usecase,
	}
}

func (h *ForgotPasswordHandler) Handle(ctx context.Context, msg *pb.Message) error {
	data := msg.GetUserForgotPasswordData()
	if data == nil {
		return fmt.Errorf("payload is  nill or not UserForgotPasswordData  for message ID %s", msg.Id)
	}

	userIDUUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	return h.usecase.Execute(ctx, userIDUUID, data.Token)

}
