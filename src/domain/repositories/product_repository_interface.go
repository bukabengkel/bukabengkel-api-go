package repositories

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/domain/entity"
)

type ProductRepositoryFilter struct {
	Name string
}

type ProductRepositoryInterface interface {
	List(ctx context.Context, page int, perPage int, sort string, filter ProductRepositoryFilter) ([]*entity.Product, int, error)
}
