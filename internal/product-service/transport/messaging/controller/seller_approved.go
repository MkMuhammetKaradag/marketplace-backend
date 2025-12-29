package controller

import (
	"context"

	"fmt"

	"marketplace/internal/product-service/transport/messaging/usecase"
	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type SellerApprovedHandler struct {
	usecase usecase.SellerApprovedUseCase
}

func NewSellerApprovedHandler(usecase usecase.SellerApprovedUseCase) *SellerApprovedHandler {
	return &SellerApprovedHandler{
		usecase: usecase,
	}
}

func (h *SellerApprovedHandler) Handle(ctx context.Context, msg *pb.Message) error {

	data := msg.GetSellerApprovedData()
	if data == nil {
		return fmt.Errorf("payload is nil or not SellerApprovedData for message ID: %s", msg.Id)
	}

	// 2. UUID doğrulaması yap
	idUUID, err := uuid.Parse(data.SellerId) // 'event' yerine doğrudan 'data' kullan
	if err != nil {
		return fmt.Errorf("invalid seller user id format: %w", err)
	}

	// 3. Usecase'e gönder
	return h.usecase.Execute(ctx, idUUID)
}
