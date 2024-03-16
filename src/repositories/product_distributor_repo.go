package repository

import (
	"context"
	"database/sql"
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
	DistributorID *uint64
	Name          *string
	Code          *string
}

func NewProductDistributorRepository(db *bun.DB, imageRepository *ImageRepository) *ProductDistributorRepository {
	return &ProductDistributorRepository{
		db:              db,
		imageRepository: imageRepository,
	}
}

func (r *ProductDistributorRepository) queryBuilder(query *bun.SelectQuery, filter ProductDistributorRepositoryFilter) *bun.SelectQuery {
	if filter.DistributorID != nil {
		query.Where("? = ?", bun.Ident("distributor_id"), filter.DistributorID)
	}

	if filter.Name != nil {
		query.Where("? ILIKE ?", bun.Ident("name"), fmt.Sprintf("%%%s%%", *filter.Name))
	}

	if filter.Code != nil {
		query.Where("? = ?", bun.Ident("code"), filter.Code)
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

func (r *ProductDistributorRepository) Save(product *models.ProductDistributor) (*models.ProductDistributor, error) {
	_, err := r.db.NewInsert().Model(product).Returning("id").Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	err = r.db.NewSelect().Model(product).WherePK().Scan(context.TODO())
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductDistributorRepository) Update(product *models.ProductDistributor) (*models.ProductDistributor, error) {
	_, err := r.db.NewUpdate().Model(product).Where("id = ?", product.ID).Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductDistributorRepository) FindOne(filter ProductDistributorRepositoryFilter) (*models.ProductDistributor, error) {
	var product models.ProductDistributor

	sl := r.db.NewSelect().Model(&product)
	sl = r.queryBuilder(sl, filter)

	err := sl.Scan(context.TODO())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &product, nil
}
