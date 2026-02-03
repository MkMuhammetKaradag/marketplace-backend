// internal/user-service/transport/http/usecase/signup.go
package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/user-service/domain"
	pb "marketplace/pkg/proto/events"

	"github.com/google/uuid"
)

type SignUpUseCase interface {
	Execute(ctx context.Context, user *domain.User) error
}
type signUpUseCase struct {
	userRepository domain.UserRepository
	messaging      domain.Messaging
}

type SignUpRequest struct {
	Username string
	Email    string
	Password string
}

func NewSignUpUseCase(repository domain.UserRepository, messaging domain.Messaging) SignUpUseCase {
	return &signUpUseCase{
		userRepository: repository,
		messaging:      messaging,
	}
}

func (u *signUpUseCase) Execute(ctx context.Context, user *domain.User) error {

	id, code, err := u.userRepository.SignUp(ctx, user)
	if err != nil {
		return err
	}
	data := &pb.UserActivationEmailData{
		ActivationId:   id.String(),
		Email:          user.Email,
		Username:       user.Username,
		ActivationCode: code,
	}
	message := &pb.Message{
		Id:          uuid.New().String(),
		Type:        pb.MessageType_USER_ACTIVATION_EMAIL,
		FromService: pb.ServiceType_USER_SERVICE,
		Critical:    true,
		RetryCount:  5,
		ToServices:  []pb.ServiceType{pb.ServiceType_NOTIFICATION_SERVICE},
		Payload:     &pb.Message_UserActivationEmailData{UserActivationEmailData: data},
	}
	if err := u.messaging.PublishMessage(ctx, message); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}
