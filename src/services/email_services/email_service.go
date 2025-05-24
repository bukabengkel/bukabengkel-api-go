package email_services

import (
	"context"
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/config"
)

type EmailService interface {
	SendEmail(ctx context.Context, to string, subject string, body string) error
}

const (
	EMAIL_SERVICE_MAILTRAP = "mailtrap"
	EMAIL_SERVICE_RESEND   = "resend"
)

func NewEmailService(config *config.Config) (EmailService, error) {
	if config.Env != "production" {
		return newMailtrapService(config), nil
	}

	switch config.EmailProvider.EmailServiceName {
	case EMAIL_SERVICE_RESEND:
		return newResendService(config), nil
	default:
		return nil, fmt.Errorf("email service %s not supported", config.EmailProvider.EmailServiceName)
	}
}
