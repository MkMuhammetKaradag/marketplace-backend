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
}

func NewForgotPasswordUseCase(emailProvider domain.EmailProvider, notificationRepository domain.NotificationRepository) ForgotPasseordUseCase {
	return &forgotPasswordUseCase{
		emailProvider:          emailProvider,
		notificationRepository: notificationRepository,
	}
}

func (u *forgotPasswordUseCase) Execute(ctx context.Context, userID uuid.UUID, token string) error {
	user, err := u.notificationRepository.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	html := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <style>
            body { font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; }
            .container { max-width: 600px; margin: 20px auto; padding: 20px; border: 1px solid #e1e1e1; border-radius: 10px; }
            .header { text-align: center; padding-bottom: 20px; border-bottom: 2px solid #f4f4f4; }
            .content { padding: 30px 0; text-align: center; }
            .seller-details { background-color: #f8f9fa; border-radius: 8px; padding: 20px; margin: 20px 0; }
            .price { font-size: 24px; color: #28a745; font-weight: bold; }
            .status-badge { background-color: #28a745; color: #fff; padding: 5px 15px; border-radius: 20px; font-size: 14px; font-weight: bold; }
            .footer { text-align: center; font-size: 12px; color: #888; border-top: 1px solid #f4f4f4; padding-top: 20px; }
            .btn { background-color: #28a745; color: white; padding: 12px 25px; text-decoration: none; border-radius: 5px; display: inline-block; margin-top: 20px; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h2>Şifre Sıfırlama </h2>
            </div>
            <div class="content">
                <p>Merhaba <strong>%s</strong>,</p>
                <p>Şifre sıfırlama isteğiniz onaylandı.</p>

                <div class="seller-details">
                     <p>Şifrenizi sıfırlamak için aşağıdaki linke tıklayın.</p>
                </div>
               
                
                <a href="http://localhost:3000/reset-password?token=%s" class="btn">Şifre Sıfırla</a>
            </div>
            <div class="footer">
                &copy; 2026 Marketplace Inc. | İstanbul, Türkiye <br>
                Bu bir otomatik bilgilendirme e-postasıdır.
            </div>
        </div>
    </body>
    </html>
    `, user.Username, token)

	return u.emailProvider.SendEmail(user.Email, "Şifre Sıfırlama", html)

}
