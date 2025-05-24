package email_services

import (
	"context"
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/config"
	resend "github.com/resend/resend-go/v2"
)

type ResendService struct {
	client *resend.Client
}

func newResendService(config *config.Config) *ResendService {
	return &ResendService{client: resend.NewClient(config.EmailProvider.EmailAPIKey)}
}

func (s *ResendService) SendEmail(ctx context.Context, to string, subject string, body string) error {
	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{"delivered@resend.dev"},
		Subject: "hello world",
		Html:    "<p>it works!</p>",
	}

	sent, err := s.client.Emails.SendWithContext(ctx, params)

	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Id)

	return nil
}
