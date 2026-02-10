// internal/notification-service/transport/messaging/controller/payment_failed.go
package controller

import (
	"context"

	"fmt"

	"marketplace/internal/notification-service/transport/messaging/usecase"
	eventsProto "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type PaymentFailedHandler struct {
	usecase usecase.PaymentFailedUseCase
}

func NewPaymentFailedHandler(usecase usecase.PaymentFailedUseCase) *PaymentFailedHandler {
	return &PaymentFailedHandler{
		usecase: usecase,
	}
}

func (h *PaymentFailedHandler) Handle(ctx context.Context, msg *eventsProto.Message) error {

	data := msg.GetPaymentFailedData()
	if data == nil {
		return fmt.Errorf("payload is nil or not PaymentFailedData for message ID: %s", msg.Id)
	}

	orderIDUUID, err := uuid.Parse(data.OrderId)
	if err != nil {
		return fmt.Errorf("invalid order id format: %w", err)
	}
	userIDUUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return fmt.Errorf("Invalid user id format:%w", err)
	}

	return h.usecase.Execute(ctx, orderIDUUID, userIDUUID)
}
