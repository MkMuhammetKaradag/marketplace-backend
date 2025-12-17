package controller

import (
	"context"
	"fmt"

	"marketplace/internal/user-service/transport/messaging/usecase"
	"marketplace/pkg/messaging"

	"github.com/google/uuid"
)

type SellerApprovedUserHandler struct {
	usecase usecase.SellerApprovedUseCase
}

func NewSellerApprovedHandler(usecase usecase.SellerApprovedUseCase) *SellerApprovedUserHandler {
	return &SellerApprovedUserHandler{
		usecase: usecase,
	}
}

func (h *SellerApprovedUserHandler) Handle(ctx context.Context, msg *messaging.Message) error {

	data := msg.GetPayload()
	if data == nil {
		return fmt.Errorf("UserCreatedData payload is nil for message ID: %s", msg.Id)
	}
	fmt.Println("Seller approved use case executed", data)
	idUUID, err := uuid.Parse(data.(map[string]interface{})["sellerUserId"].(string))
	if err != nil {
		return fmt.Errorf("UserCreatedData payload is nil for message ID: %s", msg.Id)
	}

	return h.usecase.Execute(ctx, idUUID)
}
