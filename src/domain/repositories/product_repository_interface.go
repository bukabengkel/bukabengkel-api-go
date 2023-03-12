package repositories

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/domain/entity"
)

type ProductRepositoryInterface interface {
	List(ctx context.Context, page int, perPage int, sort string, filter entity.ProductEntityRepositoryFilter) ([]*entity.ProductEntity, int, error)
}
