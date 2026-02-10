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
}

func NewUserActivationUseCase(emailProvider domain.EmailProvider) UserActivationUseCase {
	return &userActivationUseCase{
		emailProvider: emailProvider,
	}
}

func (u *userActivationUseCase) Execute(ctx context.Context, activationID uuid.UUID, userEmail string, userName string, userActivationCode string) error {

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
            .code-container { background-color: #f8f9fa; border: 2px dashed #007bff; border-radius: 8px; padding: 20px; margin: 20px 0; font-size: 32px; font-weight: bold; letter-spacing: 5px; color: #007bff; }
            .footer { text-align: center; font-size: 12px; color: #888; border-top: 1px solid #f4f4f4; padding-top: 20px; }
            .btn { background-color: #007bff; color: white; padding: 12px 25px; text-decoration: none; border-radius: 5px; display: inline-block; margin-top: 20px; }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <h2>Marketplace Aktivasyonu</h2>
            </div>
            <div class="content">
                <p>Merhaba <strong>%s</strong>,</p>
                <p>Marketplace ailesine hoş geldin! Hesabını aktifleştirmek ve alışverişe başlamak için aşağıdaki onay kodunu kullanabilirsin:</p>
                <a href="http://localhost:3000/activate?activation_id=%s">Aktivasyon Kodu</a>
                <div class="code-container">
                    %s
                </div>
                
                <p>Bu kod 30 dakika boyunca geçerlidir. Eğer bu isteği sen yapmadıysan bu e-postayı görmezden gelebilirsin.</p>
            </div>
            <div class="footer">
                &copy; 2026 Marketplace Inc. | İstanbul, Türkiye <br>
                Sana daha iyi hizmet verebilmek için buradayız.
            </div>
        </div>
    </body>
    </html>
    `, userName, activationID, userActivationCode)

	// Not: Konu kısmını "Sipariş Onayı"ndan "Hesap Aktivasyon Kodu" olarak güncelledim.
	err := u.emailProvider.SendEmail(userEmail, "Hesap Aktivasyon Kodu", html)
	if err != nil {
		return fmt.Errorf("failed to send activation email: %w", err)
	}

	return nil

}
