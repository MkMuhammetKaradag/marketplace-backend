package controller

import (
	"context"

	"fmt"

	"marketplace/internal/product-service/transport/messaging/usecase"
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

	// 2. UUID doğrulaması yap
	orderIDUUID, err := uuid.Parse(data.OrderId) // 'event' yerine doğrudan 'data' kullan
	if err != nil {
		return fmt.Errorf("invalid seller user id format: %w", err)
	}
	userIDUUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return fmt.Errorf("invalid seller user id format: %w", err)
	}

	// 3. Usecase'e gönder
	return h.usecase.Execute(ctx, orderIDUUID, userIDUUID, data.Items)
}
