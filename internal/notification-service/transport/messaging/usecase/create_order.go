// internal/notification-service/transport/messaging/usecase/create_order.go
package usecase

import (
	"context"
	"fmt"
	"marketplace/internal/notification-service/domain"

	"github.com/google/uuid"
)

type OrderCreatedUseCase interface {
	Execute(ctx context.Context, userID, orderID uuid.UUID, totalPrice float64) error
}
type orderCreatedUseCase struct {
	emailProvider          domain.EmailProvider
	notificationRepository domain.NotificationRepository
	templateMgr            domain.TemplateManager
}

func NewOrderCreatedUseCase(emailProvider domain.EmailProvider, notificationRepository domain.NotificationRepository, templateMgr domain.TemplateManager) OrderCreatedUseCase {
	return &orderCreatedUseCase{
		emailProvider:          emailProvider,
		notificationRepository: notificationRepository,
		templateMgr:            templateMgr,
	}
}

func (u *orderCreatedUseCase) Execute(ctx context.Context, userID, orderID uuid.UUID, totalPrice float64) error {

	user, err := u.notificationRepository.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"Username":   user.Username,
		"OrderID":    orderID,
		"TotalPrice": totalPrice,
	}

	html, err := u.templateMgr.Render("create_order.html", data)
	if err != nil {
		return err
	}

	err = u.emailProvider.SendEmail(user.Email, "Create Order", html)
	if err != nil {
		return fmt.Errorf("failed to send create order email: %w", err)
	}

	return nil

}
