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
}

func NewPaymentFailedUseCase(repository domain.NotificationRepository, provider domain.EmailProvider) PaymentFailedUseCase {
	return &paymentFailedUseCase{
		repository:    repository,
		emailProvider: provider,
	}
}

func (u *paymentFailedUseCase) Execute(ctx context.Context, orderID, userID uuid.UUID) error {

	user, err := u.repository.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	html := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <style>
            .container { max-width: 600px; margin: 0 auto; font-family: sans-serif; border: 1px solid #e1e1e1; border-radius: 8px; overflow: hidden; }
            .header { background-color: #a72828ff; color: white; text-align: center; padding: 20px; }
            .content { padding: 30px; text-align: center; line-height: 1.6; }
            .order-box { background-color: #f9f9f9; border: 1px solid #ddd; padding: 15px; margin: 20px 0; border-radius: 5px; }
            .amount { font-size: 24px; font-weight: bold; color: #a72828ff; }
            .footer { background-color: #f4f4f4; padding: 15px; text-align: center; font-size: 12px; color: #777; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>Ödemeniz Onaylanmadı!❗</h1>
            </div>
            <div class="content">
                <p>Merhaba <strong>%s</strong>,</p>
                <p>Harika haber! <strong>#%s</strong> numaralı siparişinizin ödemesi başarıyla alınamadı.</p>
            </div>
            <div class="footer">
                Marketplace Inc. | Güvenli Alışverişin Adresi
            </div>
        </div>
    </body>
    </html>
    `, user.Username, orderID.String())

	err = u.emailProvider.SendEmail(user.Email, "Create Order", html)
	if err != nil {
		return fmt.Errorf("failed to send create order email: %w", err)
	}

	return nil
}
