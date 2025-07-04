package usecase

import (
	"context"

	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	distributor_services "github.com/peang/bukabengkel-api-go/src/services/distributor_services"
	"github.com/peang/bukabengkel-api-go/src/services/email_services"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type webhookUsecase struct {
	logger               utils.Logger
	orderDistributorRepo *repository.OrderDistributorRepository
	userStoreRepo        *repository.UserStoreAggregateRepository
	emailService         email_services.EmailService
	distributorService   distributor_services.AsianAccessoriesService
}

type WebhookUsecase interface {
	MidtransWebhook(ctx context.Context, dto request.MidtransWebhookDTO)
}

func NewWebhookUsecase(logger utils.Logger, orderDistributorRepo *repository.OrderDistributorRepository, userStoreRepo *repository.UserStoreAggregateRepository, emailService email_services.EmailService, distributorService distributor_services.AsianAccessoriesService) WebhookUsecase {
	return &webhookUsecase{
		logger:               logger,
		orderDistributorRepo: orderDistributorRepo,
		userStoreRepo:        userStoreRepo,
		emailService:         emailService,
		distributorService:   distributorService,
	}
}
