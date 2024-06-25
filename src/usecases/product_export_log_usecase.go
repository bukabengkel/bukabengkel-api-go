package usecase

import (
	"context"
	"strconv"

	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type ProductExportLogUsecase interface {
	List(ctx context.Context, dto request.ProductExportLogListDTO) (*[]models.ProductExportLog, int, error)
}

type productExportLogUsecase struct {
	productExportLogRepository repository.ProductExportLogRepository
}

func NewProductExportLogUsecase(
	productExportLogRepository *repository.ProductExportLogRepository,
) ProductExportLogUsecase {
	return &productExportLogUsecase{
		productExportLogRepository: *productExportLogRepository,
	}
}

func (u *productExportLogUsecase) List(ctx context.Context, dto request.ProductExportLogListDTO) (*[]models.ProductExportLog, int, error) {
	filter := repository.ProductExportLogRepositoryFilter{
		StoreID: &dto.StoreID,
		UserID:  &dto.UserID,
	}

	page, err := strconv.Atoi(dto.Page)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(dto.PerPage)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	sort := "created_at"
	if dto.Sort != "" {
		sort = dto.Sort
	}

	products, count, err := u.productExportLogRepository.List(ctx, page, perPage, sort, filter)
	if err != nil {
		err = utils.NewInternalServerError(err)
	}

	return products, count, nil
}
