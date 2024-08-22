package repository

import (
	"context"
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/peang/bukabengkel-api-go/src/utils"
	"github.com/uptrace/bun"
)

type DistributorRepository struct {
	db *bun.DB
}

type DistributorRepositoryFilter struct {
	Name *string
}

type DistributorRepositoryValues struct {
	RemoteUpdate *bool
}

func NewDistributorRepository(db *bun.DB) *DistributorRepository {
	return &DistributorRepository{
		db: db,
	}
}

func (r *DistributorRepository) queryBuilder(query *bun.SelectQuery, filter DistributorRepositoryFilter) *bun.SelectQuery {
	if filter.Name != nil {
		query.Where("? ILIKE ?", bun.Ident("distributor.name"), fmt.Sprintf("%%%s%%", *filter.Name))
	}

	return query
}

func (r *DistributorRepository) List(ctx context.Context, page int, perPage int, sort string, filter DistributorRepositoryFilter) (*[]models.Distributor, int, error) {
	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimit(page, perPage)

	var distributors []models.Distributor
	sl := r.db.NewSelect().Model(&distributors)
	sl = r.queryBuilder(sl, filter)

	count, err := sl.
		Limit(limit).
		Offset(offset).
		OrderExpr(sorts).
		ScanAndCount(context.TODO())
	if err != nil {
		return nil, 0, err
	}

	if len(distributors) == 0 {
		return &[]models.Distributor{}, count, nil
	}

	return &distributors, count, nil
}
