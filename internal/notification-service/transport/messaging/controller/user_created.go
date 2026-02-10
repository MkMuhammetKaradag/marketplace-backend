// internal/notification-service/transport/messaging/controller/user_created.go
package controller

import (
	"context"

	"fmt"

	"marketplace/internal/notification-service/transport/messaging/usecase"
	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type UserCreatedHandler struct {
	usecase usecase.UserCreatedUseCase
}

func NewUserCreatedHandler(usecase usecase.UserCreatedUseCase) *UserCreatedHandler {
	return &UserCreatedHandler{
		usecase: usecase,
	}
}

func (h *UserCreatedHandler) Handle(ctx context.Context, msg *pb.Message) error {

	data := msg.GetUserCreatedData()
	if data == nil {
		return fmt.Errorf("payload is nil or not UserCreatedData for message ID: %s", msg.Id)
	}

	userIDUUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	return h.usecase.Execute(ctx, userIDUUID, data.Username, data.Email)
}
