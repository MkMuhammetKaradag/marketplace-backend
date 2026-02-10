// internal/notification-service/transport/messaging/usecase/reject_seller.go
package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/notification-service/domain"

	"github.com/google/uuid"
)

type RejectSellerUseCase interface {
	Execute(ctx context.Context, userID uuid.UUID, reason string) error
}
type rejectSellerUseCase struct {
	emailProvider          domain.EmailProvider
	notificationRepository domain.NotificationRepository
	templateMgr            domain.TemplateManager
}

func NewRejectSellerUseCase(emailProvider domain.EmailProvider, notificationRepository domain.NotificationRepository, templateMgr domain.TemplateManager) RejectSellerUseCase {
	return &rejectSellerUseCase{
		emailProvider:          emailProvider,
		notificationRepository: notificationRepository,
		templateMgr:            templateMgr,
	}
}

func (u *rejectSellerUseCase) Execute(ctx context.Context, userID uuid.UUID, reason string) error {

	user, err := u.notificationRepository.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"Username": user.Username,
		"Reason":   reason,
	}

	html, err := u.templateMgr.Render("reject_seller.html", data)
	if err != nil {
		return err
	}

	err = u.emailProvider.SendEmail(user.Email, "Reject Seller", html)
	if err != nil {
		return fmt.Errorf("failed to send reject seller email: %w", err)
	}

	return nil

}
