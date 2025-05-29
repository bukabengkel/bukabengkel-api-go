package usecase

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/transport/request"
)

type webhookUsecase struct{}

type WebhookUsecase interface {
	MidtransWebhook(ctx context.Context, dto request.MidtransWebhookDTO) error
}

func NewWebhookUsecase() WebhookUsecase {
	return &webhookUsecase{}
}
