package controller

import (
	"context"

	"fmt"

	"marketplace/internal/basket-service/transport/messaging/usecase"
	eventsProto "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type PaymentSuccessHandler struct {
	usecase usecase.PaymentSuccessUseCase
}

func NewPaymentSuccessHandler(usecase usecase.PaymentSuccessUseCase) *PaymentSuccessHandler {
	return &PaymentSuccessHandler{
		usecase: usecase,
	}
}

func (h *PaymentSuccessHandler) Handle(ctx context.Context, msg *eventsProto.Message) error {

	data := msg.GetPaymentSuccessfulData()
	if data == nil {
		return fmt.Errorf("payload is nil or not PaymentSuccessfulData for message ID: %s", msg.Id)
	}

	userIDUUID, err := uuid.Parse(data.UserId)
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	return h.usecase.Execute(ctx, userIDUUID)
}
