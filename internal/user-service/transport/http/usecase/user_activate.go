package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"
	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type UserActivateUseCase interface {
	Execute(ctx context.Context, activationID uuid.UUID, activationCode string) error
}
type UseractivateUseCase struct {
	userRepository domain.UserRepository
	messaging      domain.Messaging
}

func NewUserActivateUseCase(repository domain.UserRepository, messaging domain.Messaging) UserActivateUseCase {
	return &UseractivateUseCase{
		userRepository: repository,
		messaging:      messaging,
	}
}

func (u *UseractivateUseCase) Execute(ctx context.Context, activationID uuid.UUID, activationCode string) error {

	user, err := u.userRepository.UserActivate(ctx, activationID, activationCode)
	if err != nil {
		return err
	}
	data := &pb.UserCreatedData{
		UserId:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}

	message := &pb.Message{
		Id:          uuid.New().String(),
		Type:        pb.MessageType_USER_CREATED,
		FromService: pb.ServiceType_USER_SERVICE,
		Critical:    true,
		RetryCount:  0,
		ToServices:  []pb.ServiceType{pb.ServiceType_SELLER_SERVICE, pb.ServiceType_PRODUCT_SERVICE, pb.ServiceType_NOTIFICATION_SERVICE},
		Payload:     &pb.Message_UserCreatedData{UserCreatedData: data},
	}
	if err := u.messaging.PublishMessage(ctx, message); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
