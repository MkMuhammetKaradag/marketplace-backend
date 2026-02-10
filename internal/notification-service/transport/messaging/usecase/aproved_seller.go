// internal/notification-service/transport/messaging/usecase/aproved_seller.go
package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/notification-service/domain"

	"github.com/google/uuid"
)

type ApproveSellerUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}
type approveSellerUseCase struct {
	emailProvider          domain.EmailProvider
	templateMgr            domain.TemplateManager
	notificationRepository domain.NotificationRepository
}

func NewApproveSellerUseCase(emailProvider domain.EmailProvider, templateMgr domain.TemplateManager, notificationRepository domain.NotificationRepository) ApproveSellerUseCase {
	return &approveSellerUseCase{
		emailProvider:          emailProvider,
		notificationRepository: notificationRepository,
		templateMgr:            templateMgr,
	}
}

func (u *approveSellerUseCase) Execute(ctx context.Context, userID uuid.UUID) error {

	user, err := u.notificationRepository.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"Username": user.Username,
	}

	html, err := u.templateMgr.Render("approve_seller.html", data)
	if err != nil {
		return err
	}
	err = u.emailProvider.SendEmail(user.Email, "Approve Seller", html)
	if err != nil {
		return fmt.Errorf("failed to send approve seller email: %w", err)
	}

	return nil

}
