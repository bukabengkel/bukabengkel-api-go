package usecase

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type ProductDistributorUsecase interface {
	List(ctx context.Context, page int, perPage int, sort string, filter repository.ProductDistributorRepositoryFilter) (*[]models.ProductDistributor, int, error)
}

type productDistributorUsecase struct {
	productDistributorRepository repository.ProductDistributorRepository
}

func NewProductDistributorUsecase(
	productDistributorRepository *repository.ProductDistributorRepository,
) *productDistributorUsecase {
	return &productDistributorUsecase{
		productDistributorRepository: *productDistributorRepository,
	}
}

func (u *productDistributorUsecase) List(ctx context.Context, page int, perPage int, sort string, filter repository.ProductDistributorRepositoryFilter) (*[]models.ProductDistributor, int, error) {
	products, count, err := u.productDistributorRepository.List(ctx, page, perPage, sort, filter)
	if err != nil {
		err = utils.NewInternalServerError(err)
	}

	return products, count, nil
}
