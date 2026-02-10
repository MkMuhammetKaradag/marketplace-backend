// internal/notification-service/transport/messaging/usecase/payment_success.go
package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/notification-service/domain"

	"github.com/google/uuid"
)

type PaymentSuccessUseCase interface {
	Execute(ctx context.Context, orderID, userID uuid.UUID, amount float64) error
}
type paymentSuccessUseCase struct {
	repository    domain.NotificationRepository
	emailProvider domain.EmailProvider
	templateMgr   domain.TemplateManager
}

func NewPaymentSuccessUseCase(repository domain.NotificationRepository, provider domain.EmailProvider, templateMgr domain.TemplateManager) PaymentSuccessUseCase {
	return &paymentSuccessUseCase{
		repository:    repository,
		emailProvider: provider,
		templateMgr:   templateMgr,
	}
}

func (u *paymentSuccessUseCase) Execute(ctx context.Context, orderID, userID uuid.UUID, amount float64) error {

	user, err := u.repository.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"Username": user.Username,
		"OrderID":  orderID.String(),
		"Amount":   amount,
	}

	html, err := u.templateMgr.Render("payment_success.html", data)
	if err != nil {
		return err
	}

	err = u.emailProvider.SendEmail(user.Email, "Payment Success", html)
	if err != nil {
		return fmt.Errorf("failed to send payment success email: %w", err)
	}

	return nil
}
