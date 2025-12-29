package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/seller-service/domain"
	pb "marketplace/pkg/proto/events"

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

	data := &pb.SellerApprovedData{
		SellerId:   sellerId,
		ApprovedBy: approvedBy,
		UserId:     sellerUserId,
	}
	message := &pb.Message{
		Id:          uuid.New().String(),
		Type:        pb.MessageType_SELLER_APPROVED,
		FromService: pb.ServiceType_SELLER_SERVICE,
		ToServices:  []pb.ServiceType{pb.ServiceType_USER_SERVICE, pb.ServiceType_PRODUCT_SERVICE},
		Payload:     &pb.Message_SellerApprovedData{SellerApprovedData: data},
	}

	if err := u.messaging.PublishMessage(ctx, message); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}
