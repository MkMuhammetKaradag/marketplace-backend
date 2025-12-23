package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/seller-service/domain"
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

	// data := &pb.SellerRejectedData{
	// 	SellerId:   sellerUserId,
	// 	RejectedBy: rejectedBy,
	// 	Reason:     reason,
	// }
	// message := &pb.Message{
	// 	Id:          uuid.New().String(),
	// 	Type:        pb.MessageType_SELLER_REJECTED,
	// 	FromService: pb.ServiceType_SELLER_SERVICE,
	// 	ToServices:  []pb.ServiceType{pb.ServiceType_USER_SERVICE},
	// 	Payload:     &pb.Message_SellerRejectedData{SellerRejectedData: data},
	// }

	// if err := u.messaging.PublishMessage(ctx, message); err != nil {
	// 	return fmt.Errorf("failed to publish message: %w", err)
	// }
	return nil
}
