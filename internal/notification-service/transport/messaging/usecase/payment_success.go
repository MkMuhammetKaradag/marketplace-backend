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
}

func NewPaymentSuccessUseCase(repository domain.NotificationRepository, provider domain.EmailProvider) PaymentSuccessUseCase {
	return &paymentSuccessUseCase{
		repository:    repository,
		emailProvider: provider,
	}
}

func (u *paymentSuccessUseCase) Execute(ctx context.Context, orderID, userID uuid.UUID, amount float64) error {

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
            .header { background-color: #28a745; color: white; text-align: center; padding: 20px; }
            .content { padding: 30px; text-align: center; line-height: 1.6; }
            .order-box { background-color: #f9f9f9; border: 1px solid #ddd; padding: 15px; margin: 20px 0; border-radius: 5px; }
            .amount { font-size: 24px; font-weight: bold; color: #28a745; }
            .footer { background-color: #f4f4f4; padding: 15px; text-align: center; font-size: 12px; color: #777; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h1>Ödemeniz Onaylandı! ✅</h1>
            </div>
            <div class="content">
                <p>Merhaba <strong>%s</strong>,</p>
                <p>Harika haber! <strong>#%s</strong> numaralı siparişinizin ödemesi başarıyla alındı.</p>
                <div class="order-box">
                    <p>Ödenen Tutar</p>
                    <div class="amount">%.2f TL</div>
                </div>
                <p>Siparişiniz hazırlık aşamasına alınmıştır. Kargoya verildiğinde sizi tekrar bilgilendireceğiz.</p>
            </div>
            <div class="footer">
                Marketplace Inc. | Güvenli Alışverişin Adresi
            </div>
        </div>
    </body>
    </html>
    `, user.Username, orderID.String(), amount)

	err = u.emailProvider.SendEmail(user.Email, "Create Order", html)
	if err != nil {
		return fmt.Errorf("failed to send create order email: %w", err)
	}

	return nil
}
