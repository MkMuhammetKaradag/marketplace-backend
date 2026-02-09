package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"
	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type ForgotPasswordUseCase interface {
	Execute(ctx context.Context, identifier string) error
}
type forgotPasswordUseCase struct {
	repo      domain.UserRepository
	messaging domain.Messaging
}

func NewForgotPasswordUseCase(repo domain.UserRepository, messaging domain.Messaging) ForgotPasswordUseCase {
	return &forgotPasswordUseCase{
		repo:      repo,
		messaging: messaging,
	}
}

func (u *forgotPasswordUseCase) Execute(ctx context.Context, identifier string) error {

	result, err := u.repo.ForgotPassword(ctx, identifier)
	if err != nil {
		return err
	}
	data := &pb.UserForgotPasswordData{
		UserId: result.UserID.String(),
		Token:  result.Token,
	}
	message := &pb.Message{
		Id:          uuid.New().String(),
		Type:        pb.MessageType_USER_FORGOT_PASSWORD,
		FromService: pb.ServiceType_USER_SERVICE,
		Critical:    true,
		RetryCount:  5,
		ToServices:  []pb.ServiceType{pb.ServiceType_NOTIFICATION_SERVICE},
		Payload:     &pb.Message_UserForgotPasswordData{UserForgotPasswordData: data},
	}
	if err := u.messaging.PublishMessage(ctx, message); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}
