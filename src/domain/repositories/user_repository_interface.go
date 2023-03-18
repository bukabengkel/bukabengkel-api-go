package repositories

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/domain/entity"
)

type UserRepositoryFilter struct {
	Name string
}

type UserRepositoryInterface interface {
	Detail(ctx context.Context, id int) (*entity.User, int, error)
	List(ctx context.Context, page int, perPage int, sort string, filter UserRepositoryFilter) ([]*entity.User, int, error)
}
