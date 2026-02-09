package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/seller-service/domain"
	pEvents "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type RejectSellerUseCase interface {
	Execute(ctx context.Context, sellerId, rejectedBy string, reason string) error
}
type rejectSellerUseCase struct {
	sellerRepository domain.SellerRepository
	messaging        domain.Messaging
}

func NewRejectSellerUseCase(repository domain.SellerRepository, messaging domain.Messaging) RejectSellerUseCase {
	return &rejectSellerUseCase{
		sellerRepository: repository,
		messaging:        messaging,
	}
}

func (u *rejectSellerUseCase) Execute(ctx context.Context, sellerId, rejectedBy string, reason string) error {

	sellerUserId, err := u.sellerRepository.RejectSeller(ctx, sellerId, rejectedBy, reason)
	if err != nil {
		return err
	}
	fmt.Println("seller user id ", sellerUserId)

	// Publish message to Kafka

	data := &pEvents.SellerRejectedData{
		SellerId:   sellerUserId,
		RejectedBy: rejectedBy,
		Reason:     reason,
	}
	message := &pEvents.Message{
		Id:          uuid.New().String(),
		Type:        pEvents.MessageType_SELLER_REJECTED,
		FromService: pEvents.ServiceType_SELLER_SERVICE,
		ToServices:  []pEvents.ServiceType{pEvents.ServiceType_NOTIFICATION_SERVICE},
		Payload:     &pEvents.Message_SellerRejectedData{SellerRejectedData: data},
	}

	if err := u.messaging.PublishMessage(ctx, message); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}
