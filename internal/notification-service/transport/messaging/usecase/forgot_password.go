// internal/notification-service/transport/messaging/usecase/forgot_password.go
package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/notification-service/domain"

	"github.com/google/uuid"
)

type ForgotPasseordUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, token string) error
}

type forgotPasswordUseCase struct {
	emailProvider          domain.EmailProvider
	notificationRepository domain.NotificationRepository
	templateMgr            domain.TemplateManager
}

func NewForgotPasswordUseCase(emailProvider domain.EmailProvider, notificationRepository domain.NotificationRepository, templateMgr domain.TemplateManager) ForgotPasseordUseCase {
	return &forgotPasswordUseCase{
		emailProvider:          emailProvider,
		notificationRepository: notificationRepository,
		templateMgr:            templateMgr,
	}
}

func (u *forgotPasswordUseCase) Execute(ctx context.Context, userID uuid.UUID, token string) error {
	user, err := u.notificationRepository.GetUser(ctx, userID)
	if err != nil {
		fmt.Println("err-user not found:", err)
		return err
	}

	data := map[string]interface{}{
		"Username": user.Username,
		"Token":    token,
	}

	html, err := u.templateMgr.Render("forgot_password.html", data)
	if err != nil {
		fmt.Println("err:", err)
		return fmt.Errorf("failed to render forgot password template: %w", err)
	}

	return u.emailProvider.SendEmail(user.Email, "Şifre Sıfırlama", html)

}
