package usecase

import (
	"context"
	"strconv"

	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

func (u *orderDistributorUsecase) List(ctx context.Context, dto request.OrderDistributorListDTO) (*[]models.OrderDistributor, bool, error) {
	filter := repository.OrderDistributorRepositoryFilter{
		StoreID: &dto.StoreID,
	}

	page, err := strconv.Atoi(dto.Page)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(dto.PerPage)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	if dto.OrderID != "" && len(dto.OrderID) >= 3 {
		filter.Key = utils.String(dto.OrderID)
	}

	sort := "created_at"
	if dto.Sort != "" {
		sort = dto.Sort
	}

	orderDistributors, next, err := u.orderDistributorRepository.List(ctx, page, perPage, sort, filter)
	if err != nil {
		err = utils.NewInternalServerError(err)
	}

	return &orderDistributors, next, nil
}
