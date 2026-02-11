// internal/notification-service/transport/messaging/usecase/user_activation.go
package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/notification-service/domain"

	"github.com/google/uuid"
)

type UserActivationUseCase interface {
	Execute(ctx context.Context, activationID uuid.UUID, userEmail string, userName string, userActivationCode string) error
}
type userActivationUseCase struct {
	emailProvider domain.EmailProvider
	templateMgr   domain.TemplateManager
}

func NewUserActivationUseCase(emailProvider domain.EmailProvider, templateMgr domain.TemplateManager) UserActivationUseCase {
	return &userActivationUseCase{
		emailProvider: emailProvider,
		templateMgr:   templateMgr,
	}
}

func (u *userActivationUseCase) Execute(ctx context.Context, activationID uuid.UUID, userEmail string, userName string, userActivationCode string) error {

	data := map[string]interface{}{
		"Username":       userName,
		"ActivationID":   activationID.String(),
		"ActivationCode": userActivationCode,
	}
	fmt.Println("data", data)

	html, err := u.templateMgr.Render("user_activation.html", data)
	if err != nil {
		return err
	}
	err = u.emailProvider.SendEmail(userEmail, "Hesap Aktivasyon Kodu", html)
	if err != nil {
		return fmt.Errorf("failed to send activation email: %w", err)
	}
	fmt.Println("email sent successfully!", userEmail)

	return nil

}
