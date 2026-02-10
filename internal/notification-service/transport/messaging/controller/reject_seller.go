// internal/notification-service/transport/messaging/controller/reject_seller.go
package controller

import (
	"context"

	"fmt"

	"marketplace/internal/notification-service/transport/messaging/usecase"
	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type RejectSellerHandler struct {
	usecase usecase.RejectSellerUseCase
}

func NewRejectSellerHandler(usecase usecase.RejectSellerUseCase) *RejectSellerHandler {
	return &RejectSellerHandler{
		usecase: usecase,
	}
}

func (h *RejectSellerHandler) Handle(ctx context.Context, msg *pb.Message) error {

	data := msg.GetSellerRejectedData()
	if data == nil {
		return fmt.Errorf("payload is nil or not OrderCreatedData for message ID: %s", msg.Id)
	}

	userIDUUID, err := uuid.Parse(data.SellerId)
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	reason := data.Reason

	return h.usecase.Execute(ctx, userIDUUID, reason)
}
