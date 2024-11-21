package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strconv"

	"github.com/gotidy/ptr"
	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type ProductUsecase interface {
	List(ctx context.Context, dto request.ProductListDTO) (*[]models.Product, int, error)
	ListV2(ctx context.Context, dto request.ProductListDTOV2) (*[]models.Product, string, error)
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

func (u *productUsecase) List(ctx context.Context, dto request.ProductListDTO) (*[]models.Product, int, error) {
	filter := repository.ProductRepositoryFilter{
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

	if dto.Name != "" && len(dto.Name) >= 3 {
		filter.Name = utils.String(dto.Name)
	}

	if dto.Keyword != "" && len(dto.Keyword) >= 3 {
		filter.Name = utils.String(dto.Keyword)
	}

	if dto.CategoryId != "" && dto.CategoryId != "0" {
		filter.CategoryId = utils.String(dto.CategoryId)
	}

	if dto.StockMoreThan != "" {
		stockMoreThan, err := strconv.ParseUint(dto.StockMoreThan, 10, 10)
		if err != nil {
			return nil, 0, err
		}

		filter.StockMoreThan = ptr.Of(uint(stockMoreThan))
	}

	if dto.Status != "" {
		status, err := strconv.ParseUint(dto.Status, 10, 10)
		if err != nil {
			return nil, 0, err
		}

		filter.Status = ptr.Of(uint(status))
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

func (u *productUsecase) ListV2(ctx context.Context, dto request.ProductListDTOV2) (*[]models.Product, string, error) {
	filter := repository.ProductRepositoryFilter{
		StoreID: &dto.StoreID,
	}

	limit, err := strconv.Atoi(dto.Limit)
	if err != nil || limit <= 0 {
		limit = 20
	}

	if dto.Name != "" && len(dto.Name) >= 3 {
		filter.Name = utils.String(dto.Name)
	}

	if dto.Keyword != "" && len(dto.Keyword) >= 3 {
		filter.Name = utils.String(dto.Keyword)
	}

	if dto.CategoryId != "" && dto.CategoryId != "0" {
		filter.CategoryId = utils.String(dto.CategoryId)
	}

	if dto.StockMoreThan != "" {
		stockMoreThan, err := strconv.ParseUint(dto.StockMoreThan, 10, 10)
		if err != nil {
			return nil, "", err
		}

		filter.StockMoreThan = ptr.Of(uint(stockMoreThan))
	}

	if dto.Status != "" {
		status, err := strconv.ParseUint(dto.Status, 10, 10)
		if err != nil {
			return nil, "", err
		}

		filter.Status = ptr.Of(uint(status))
	}

	var currCursor = &repository.ProductCursor{}
	if dto.Cursor != "" {
		cursorBytes, err := base64.URLEncoding.DecodeString(dto.Cursor)
		if err != nil {
			return nil, "", err
		}

		if err := json.Unmarshal(cursorBytes, &currCursor); err != nil {
			return nil, "", err
		}

	}

	sort := "name"
	if dto.Sort != "" {
		sort = dto.Sort
	}

	products, nextCursor, err := u.productRepository.ListV2(ctx, currCursor, limit, sort, filter)

	if err != nil {
		err = utils.NewInternalServerError(err)
	}

	var nextCursorStr string
	if nextCursor != nil {
		cursorBytes, _ := json.Marshal(nextCursor)
		nextCursorStr = base64.URLEncoding.EncodeToString(cursorBytes)
	}

	return products, nextCursorStr, nil
}
