package controller

import (
	"context"

	"fmt"

	"marketplace/internal/user-service/transport/messaging/usecase"
	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type SellerApprovedUserHandler struct {
	usecase usecase.SellerApprovedUseCase
}

func NewSellerApprovedHandler(usecase usecase.SellerApprovedUseCase) *SellerApprovedUserHandler {
	return &SellerApprovedUserHandler{
		usecase: usecase,
	}
}

func (h *SellerApprovedUserHandler) Handle(ctx context.Context, msg *pb.Message) error {

	data := msg.GetSellerApprovedData()
	if data == nil {
		return fmt.Errorf("payload is nil for message ID: %s", msg.Id)
	}
	fmt.Println("Seller approved use case executed", data)

	// Safely convert payload to struct
	var event pb.SellerApprovedData

	// Convert interface{} (map) back to bytes then to struct
	// This handles both map[string]interface{} (from JSON) and struct (if passed internally) cases reliably
	payloadBytes, err := proto.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	if err := proto.Unmarshal(payloadBytes, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payload to SellerApprovedData: %w", err)
	}

	idUUID, err := uuid.Parse(event.SellerId)
	if err != nil {
		return fmt.Errorf("invalid seller user id format: %w", err)
	}

	return h.usecase.Execute(ctx, idUUID)
}
