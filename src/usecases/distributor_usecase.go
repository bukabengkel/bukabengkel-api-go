package usecase

import (
	"context"
	"strconv"

	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type DistributorUsecase interface {
	List(ctx context.Context, dto request.DistributorListDTO) (*[]models.Distributor, int, error)
}

type distributorUsecase struct {
	distributorRepository repository.DistributorRepository
}

func NewDistributorUsecase(
	distributorRepository *repository.DistributorRepository,
) DistributorUsecase {
	return &distributorUsecase{
		distributorRepository: *distributorRepository,
	}
}

func (u *distributorUsecase) List(ctx context.Context, dto request.DistributorListDTO) (*[]models.Distributor, int, error) {
	filter := repository.DistributorRepositoryFilter{}

	page, err := strconv.Atoi(dto.Page)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(dto.PerPage)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	if dto.Name != "" && len(dto.Name) >= 3 {
		filter.Name = utils.String(dto.Name)
	}

	sort := "name"
	if dto.Sort != "" {
		sort = dto.Sort
	}

	products, count, err := u.distributorRepository.List(ctx, page, perPage, sort, filter)
	if err != nil {
		err = utils.NewInternalServerError(err)
	}

	return products, count, nil
}
