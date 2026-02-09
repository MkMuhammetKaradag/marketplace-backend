package controller

import (
	"context"

	"fmt"

	"marketplace/internal/notification-service/transport/messaging/usecase"
	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type ApproveSellerHandler struct {
	usecase usecase.ApproveSellerUseCase
}

func NewApproveSellerHandler(usecase usecase.ApproveSellerUseCase) *ApproveSellerHandler {
	return &ApproveSellerHandler{
		usecase: usecase,
	}
}

func (h *ApproveSellerHandler) Handle(ctx context.Context, msg *pb.Message) error {

	data := msg.GetSellerApprovedData()
	if data == nil {
		return fmt.Errorf("payload is nil or not OrderCreatedData for message ID: %s", msg.Id)
	}

	userIDUUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	return h.usecase.Execute(ctx, userIDUUID)
}
