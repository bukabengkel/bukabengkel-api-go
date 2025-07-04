package email_services

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/config"
	resend "github.com/resend/resend-go/v2"
)

type ResendService struct {
	client *resend.Client
}

func newResendService(config *config.Config) *ResendService {
	return &ResendService{client: resend.NewClient(config.EmailProvider.EmailAPIKey)}
}

func (s *ResendService) SendWaitingForPaymentEmail(ctx context.Context, to string, data WaitingForPaymentData) error {

	return nil
}
