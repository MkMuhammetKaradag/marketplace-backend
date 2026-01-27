package controller

import (
	"context"

	"fmt"

	"marketplace/internal/product-service/transport/messaging/usecase"
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

	orderIDUUID, err := uuid.Parse(data.OrderId)
	if err != nil {
		return fmt.Errorf("invalid order id format: %w", err)
	}

	return h.usecase.Execute(ctx, orderIDUUID)
}
