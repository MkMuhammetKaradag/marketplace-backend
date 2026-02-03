package domain

type EmailProvider interface {
	SendEmail(to string, subject string, htmlContent string) error
}
