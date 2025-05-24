package usecase

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

func (u *orderDistributorUsecase) Detail(ctx context.Context, dto request.OrderDistributorDetailDTO) (*models.OrderDistributor, error) {
	filter := repository.OrderDistributorRepositoryFilter{
		StoreID: &dto.StoreID,
	}

	orderDistributor, err := u.orderDistributorRepository.FindOne(ctx, filter)
	if err != nil {
		err = utils.NewInternalServerError(err)
	}

	return orderDistributor, nil
}

