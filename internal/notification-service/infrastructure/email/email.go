package email

import (
	"fmt"

	"marketplace/internal/notification-service/domain"

	"github.com/resend/resend-go/v2"
)

type resendProvider struct {
	ApiKey string
	Client *resend.Client
}

func NewResendProvider(apiKey string) domain.EmailProvider {
	client := resend.NewClient(apiKey)
	return &resendProvider{
		ApiKey: apiKey,
		Client: client,
	}
}

func (r *resendProvider) SendEmail(to string, subject string, htmlContent string) error {
	params := &resend.SendEmailRequest{
		From:    "Marketplace <onboarding@resend.dev>",
		To:      []string{"onboarding@resend.dev"}, //to
		Subject: subject,
		Html:    htmlContent,
	}

	sent, err := r.Client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("resend error: %w", err)
	}

	fmt.Printf("📧 Email sent successfully! ID: %s\n", sent.Id)
	return nil
}
