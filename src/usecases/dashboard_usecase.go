package usecase

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/models"
	repository "github.com/peang/bukabengkel-api-go/src/repositories"
	"github.com/peang/bukabengkel-api-go/src/transport/request"
)

type ReportUsecase interface {
	List(ctx context.Context, dto request.ProductListDTO) (*[]models.Product, int, error)
}

type reportUsecase struct {
	productRepository repository.ProductRepository
}

func NewReportUsecase(
	productRepository *repository.ProductRepository,
) ProductUsecase {
	return &productUsecase{
		productRepository: *productRepository,
	}
}
