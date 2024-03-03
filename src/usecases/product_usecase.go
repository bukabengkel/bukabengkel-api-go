package usecase

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/domain/entity"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

// AuthUsecase represent the todos usecase contract
type ProductUsecase interface {
	List(ctx context.Context, page int, perPage int, sort string, filter repository.ProductRepositoryFilter) (*[]entity.Product, int, error)
}

type productUsecase struct {
	productRepository repository.ProductRepository
}

func NewProductUsecase(
	productRepository *repository.ProductRepository,
) ProductUsecase {
	return &productUsecase{
		productRepository: *productRepository,
	}
}

func (u *productUsecase) List(ctx context.Context, page int, perPage int, sort string, filter repository.ProductRepositoryFilter) (*[]entity.Product, int, error) {
	products, count, err := u.productRepository.List(ctx, page, perPage, sort, filter)
	if err != nil {
		err = utils.NewInternalServerError(err)
	}

	return products, count, nil
}
