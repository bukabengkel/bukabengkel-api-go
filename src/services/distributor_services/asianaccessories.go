package distributor_services

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/config"
	"github.com/peang/bukabengkel-api-go/src/models"
)

type AsianAccessoriesService struct {
	APIKey string
}

func NewAsianAccessoriesService(config *config.Config) *AsianAccessoriesService {
	return &AsianAccessoriesService{
		APIKey: config.AsianAccessoriesAPIKey,
	}
}

func Checkout(ctx context.Context, ) (*models.OrderDistributor, error) {
	return nil, nil
}