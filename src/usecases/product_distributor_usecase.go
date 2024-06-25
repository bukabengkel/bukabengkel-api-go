package usecase

import (
	"context"
	"strconv"

	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type ProductDistributorUsecase interface {
	List(ctx context.Context, dto request.ProductDistributorListDTO) (*[]models.ProductDistributor, int, error)
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

func (u *productDistributorUsecase) List(ctx context.Context, dto request.ProductDistributorListDTO) (*[]models.ProductDistributor, int, error) {
	filter := repository.ProductDistributorRepositoryFilter{}

	page, err := strconv.Atoi(dto.Page)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(dto.PerPage)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	if dto.Keyword != "" && len(dto.Keyword) >= 3 {
		filter.Name = utils.String(dto.Keyword)

	}

	sort := "-id"
	if dto.Sort != "" {
		sort = dto.Sort
	}

	products, count, err := u.productDistributorRepository.List(ctx, page, perPage, sort, filter)
	if err != nil {
		err = utils.NewInternalServerError(err)
	}

	return products, count, nil
}
