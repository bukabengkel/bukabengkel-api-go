package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/gotidy/ptr"
	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/services"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type ProductUsecase interface {
	List(ctx context.Context, dto request.ProductListDTO) (*[]models.Product, int, error)
	Export(ctx context.Context, dto request.ProductExportDTO) (bool, error)
}

type productUsecase struct {
	productRepository          *repository.ProductRepository
	productExportLogRepository *repository.ProductExportLogRepository
}

func NewProductUsecase(
	productRepository *repository.ProductRepository,
	productExportLogRepository *repository.ProductExportLogRepository,
) ProductUsecase {
	return &productUsecase{
		productRepository:          productRepository,
		productExportLogRepository: productExportLogRepository,
	}
}

func (u *productUsecase) List(ctx context.Context, dto request.ProductListDTO) (*[]models.Product, int, error) {
	filter := repository.ProductRepositoryFilter{
		StoreID: &dto.StoreID,
	}

	page, err := strconv.Atoi(dto.Page)
	if err != nil {
		return nil, 0, err
	}
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(dto.PerPage)
	if err != nil {
		return nil, 0, err
	}
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 10
	}

	if dto.Keyword != "" && len(dto.Keyword) >= 3 {
		filter.Name = utils.String(dto.Keyword)
	}

	if dto.CategoryId != "" && dto.CategoryId != "0" {
		categoryId, err := strconv.Atoi(dto.CategoryId)
		if err != nil {
			return nil, 0, err
		}
		filter.CategoryId = utils.Uint64(categoryId)
	}

	if dto.StockMoreThan != "" && dto.StockMoreThan != "0" {
		stockMoreThan, err := strconv.ParseUint(dto.StockMoreThan, 10, 10)
		if err != nil {
			return nil, 0, err
		}

		filter.StockMoreThan = ptr.Of(uint(stockMoreThan))
	}

	sort := "name"
	if dto.Sort != "" {
		sort = dto.Sort
	}

	products, count, err := u.productRepository.List(ctx, page, perPage, sort, filter)
	if err != nil {
		err = utils.NewInternalServerError(err)
	}

	return products, count, nil
}

func (u *productUsecase) Export(ctx context.Context, dto request.ProductExportDTO) (bool, error) {
	productLog := models.ProductExportLog{
		StoreID:   dto.StoreID,
		UserID:    dto.UserID,
		Status:    models.LogDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := u.productExportLogRepository.Save(&productLog)
	if err != nil {
		return false, utils.NewInternalServerError(err)
	}

	go services.ExportProduct(u.productRepository, u.productExportLogRepository, productLog.ID, dto.StoreID, dto.CategoryId)
	return false, nil
}
