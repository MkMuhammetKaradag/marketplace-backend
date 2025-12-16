package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/seller-service/domain"
	"marketplace/pkg/messaging"

	"github.com/google/uuid"
)

type ApproveSellerUseCase interface {
	Execute(ctx context.Context, sellerId, approvedBy string) error
}
type approveSellerUseCase struct {
	sellerRepository domain.SellerRepository
	messaging        domain.Messaging
}

func NewApproveSellerUseCase(repository domain.SellerRepository, messaging domain.Messaging) ApproveSellerUseCase {
	return &approveSellerUseCase{
		sellerRepository: repository,
		messaging:        messaging,
	}
}

func (u *approveSellerUseCase) Execute(ctx context.Context, sellerId, approvedBy string) error {

	sellerUserId, err := u.sellerRepository.ApproveSeller(ctx, sellerId, approvedBy)
	if err != nil {
		return err
	}
	fmt.Println("seller user id ", sellerUserId)

	// Publish message to Kafka
	message := &messaging.Message{
		Id:          uuid.New().String(),
		Type:        messaging.MessageType_SELLER_APPROVED,
		FromService: messaging.ServiceType_SELLER_SERVICE,
		ToServices:  []messaging.ServiceType{messaging.ServiceType_USER_SERVICE},
		Payload: map[string]interface{}{
			"sellerUserId": sellerUserId,
			"approvedBy":   approvedBy,
		},
	}

	if err := u.messaging.PublishMessage(ctx, message); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}
