// internal/notification-service/transport/messaging/usecase/payment_failed.go
package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/notification-service/domain"

	"github.com/google/uuid"
)

type PaymentFailedUseCase interface {
	Execute(ctx context.Context, orderID, userID uuid.UUID) error
}
type paymentFailedUseCase struct {
	repository    domain.NotificationRepository
	emailProvider domain.EmailProvider
	templateMgr   domain.TemplateManager
}

func NewPaymentFailedUseCase(repository domain.NotificationRepository, provider domain.EmailProvider, templateMgr domain.TemplateManager) PaymentFailedUseCase {
	return &paymentFailedUseCase{
		repository:    repository,
		emailProvider: provider,
		templateMgr:   templateMgr,
	}
}

func (u *paymentFailedUseCase) Execute(ctx context.Context, orderID, userID uuid.UUID) error {

	user, err := u.repository.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"Username": user.Username,
		"OrderID":  orderID.String(),
	}

	html, err := u.templateMgr.Render("payment_failed.html", data)
	if err != nil {
		return err
	}

	err = u.emailProvider.SendEmail(user.Email, "Payment Failed", html)
	if err != nil {
		return fmt.Errorf("failed to send payment failed email: %w", err)
	}

	return nil
}
