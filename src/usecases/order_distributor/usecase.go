package usecase

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
)

type orderDistributorUsecase struct {
	orderDistributorRepository repository.OrderDistributorRepository
}
type OrderDistributorUsecase interface {
	List(ctx context.Context, dto request.OrderDistributorListDTO) (*[]models.OrderDistributor, bool, error)
	Detail(ctx context.Context, dto request.OrderDistributorDetailDTO) (*models.OrderDistributor, error)
}

func NewOrderDistributorUsecase(
	orderDistributorRepository *repository.OrderDistributorRepository,
) OrderDistributorUsecase {
	return &orderDistributorUsecase{
		orderDistributorRepository: *orderDistributorRepository,
	}
}
