// internal/notification-service/transport/messaging/controller/order_created.go
package controller

import (
	"context"

	"fmt"

	"marketplace/internal/notification-service/transport/messaging/usecase"
	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type OrderCreatedHandler struct {
	usecase usecase.OrderCreatedUseCase
}

func NewOrderCreatedHandler(usecase usecase.OrderCreatedUseCase) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		usecase: usecase,
	}
}

func (h *OrderCreatedHandler) Handle(ctx context.Context, msg *pb.Message) error {

	data := msg.GetOrderCreatedData()
	if data == nil {
		return fmt.Errorf("payload is nil or not OrderCreatedData for message ID: %s", msg.Id)
	}

	userIDUUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	orderIDUUID, err := uuid.Parse(data.OrderId)
	if err != nil {
		return fmt.Errorf("invalid order id format: %w", err)
	}
	totalPrice := data.TotalPrice

	return h.usecase.Execute(ctx, userIDUUID, orderIDUUID, totalPrice)
}
