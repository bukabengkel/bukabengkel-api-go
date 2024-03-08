package repository

import (
	"context"
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/peang/bukabengkel-api-go/src/utils"
	"github.com/uptrace/bun"
)

type ProductDistributorRepository struct {
	db              *bun.DB
	imageRepository *ImageRepository
}

type ProductDistributorRepositoryFilter struct {
	Name *string
}

func NewProductDistributorRepository(db *bun.DB, imageRepository *ImageRepository) *ProductDistributorRepository {
	return &ProductDistributorRepository{
		db:              db,
		imageRepository: imageRepository,
	}
}

func (r *ProductDistributorRepository) queryBuilder(query *bun.SelectQuery, cond ProductDistributorRepositoryFilter) *bun.SelectQuery {
	if cond.Name != nil {
		query.Where("? ILIKE ?", bun.Ident("product_distributor.name"), fmt.Sprintf("%%%s%%", *cond.Name))
	}

	return query
}

func (r *ProductDistributorRepository) List(ctx context.Context, page int, perPage int, sort string, filter ProductDistributorRepositoryFilter) (*[]models.ProductDistributor, int, error) {
	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimit(page, perPage)

	var products []models.ProductDistributor
	sl := r.db.NewSelect().Model(&products)
	sl = r.queryBuilder(sl, filter)

	count, err := sl.
		Relation("Location").
		Limit(limit).Offset(offset).OrderExpr(sorts).ScanAndCount(context.TODO())
	if err != nil {
		return nil, 0, err
	}

	if len(products) == 0 {
		return &[]models.ProductDistributor{}, count, nil
	}

	var entityProducts []models.ProductDistributor
	entityProducts = append(entityProducts, products...)

	return &entityProducts, count, nil
}
