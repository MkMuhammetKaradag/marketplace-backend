package controller

import (
	"context"

	"fmt"

	"marketplace/internal/order-service/transport/messaging/usecase"
	eventsProto "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type PaymentFailureHandler struct {
	usecase usecase.PaymentFailureUseCase
}

func NewPaymentFailureHandler(usecase usecase.PaymentFailureUseCase) *PaymentFailureHandler {
	return &PaymentFailureHandler{
		usecase: usecase,
	}
}

func (h *PaymentFailureHandler) Handle(ctx context.Context, msg *eventsProto.Message) error {

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
		return fmt.Errorf("invalid seller user id format: %w", err)
	}
	errorMessage := data.ErrorMessage
	if errorMessage == "" {
		return fmt.Errorf("invalid error message format: %w", err)
	}

	return h.usecase.Execute(ctx, orderIDUUID, userIDUUID, errorMessage)
}
