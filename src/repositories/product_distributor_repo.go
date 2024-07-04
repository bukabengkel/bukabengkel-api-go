package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/models"
	file_service "github.com/peang/bukabengkel-api-go/src/services/file_services"
	"github.com/peang/bukabengkel-api-go/src/utils"
	"github.com/uptrace/bun"
)

type ProductDistributorRepository struct {
	db          *bun.DB
	fileService file_service.FileService
}

type ProductDistributorRepositoryFilter struct {
	DistributorID *uint64
	Name          *string
	Code          *string
	RemoteUpdate  *bool
}

type ProductDistributorRepositoryValues struct {
	RemoteUpdate *bool
}

func NewProductDistributorRepository(db *bun.DB, fileService file_service.FileService) *ProductDistributorRepository {
	return &ProductDistributorRepository{
		db:          db,
		fileService: fileService,
	}
}

func (r *ProductDistributorRepository) queryBuilder(query *bun.SelectQuery, filter ProductDistributorRepositoryFilter) *bun.SelectQuery {
	if filter.DistributorID != nil {
		query.Where("? = ?", bun.Ident("product_distributor.distributor_id"), filter.DistributorID)
	}

	if filter.Name != nil {
		query.Where("? ILIKE ?", bun.Ident("product_distributor.name"), fmt.Sprintf("%%%s%%", *filter.Name))
	}

	if filter.Code != nil {
		query.Where("? = ?", bun.Ident("product_distributor.code"), filter.Code)
	}

	if filter.RemoteUpdate != nil {
		query.Where("? = ?", bun.Ident("product_distributor.remote_update"), filter.RemoteUpdate)
	}

	return query
}

func (r *ProductDistributorRepository) updateQueryBuilder(
	query *bun.UpdateQuery,
	filter ProductDistributorRepositoryFilter,
	values ProductDistributorRepositoryValues,
) *bun.UpdateQuery {
	if filter.DistributorID != nil {
		query.Where("? = ?", bun.Ident("distributor_id"), filter.DistributorID)
	}

	if values.RemoteUpdate != nil {
		query.Set("? = ?", bun.Ident("remote_update"), values.RemoteUpdate)
	}

	return query
}

func (r *ProductDistributorRepository) deleteQueryBuilder(
	query *bun.DeleteQuery,
	filter ProductDistributorRepositoryFilter,
) *bun.DeleteQuery {
	if filter.DistributorID != nil {
		query.Where("? = ?", bun.Ident("distributor_id"), filter.DistributorID)
	}

	if filter.RemoteUpdate != nil {
		query.Where("? = ?", bun.Ident("remote_update"), filter.RemoteUpdate)
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
		Relation("Distributor").
		Relation("Category").
		Limit(limit).
		Offset(offset).
		OrderExpr(sorts).
		ScanAndCount(context.TODO())
	if err != nil {
		return nil, 0, err
	}

	if len(products) == 0 {
		return &[]models.ProductDistributor{}, count, nil
	}

	var entityProducts []models.ProductDistributor
	for _, p := range products {
		p.Thumbnail = r.fileService.BuildUrl(p.Thumbnail, 200, 200)
		entityProducts = append(entityProducts, p)
	}
	// entityProducts = append(entityProducts, products...)

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

func (r *ProductDistributorRepository) UpdateWithCondition(filter ProductDistributorRepositoryFilter, values ProductDistributorRepositoryValues) (int64, error) {
	var product models.ProductDistributor

	sl := r.db.NewUpdate().Model(&product)
	sl = r.updateQueryBuilder(sl, filter, values)

	res, err := sl.Exec(context.TODO())
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
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

func (r *ProductDistributorRepository) Delete(product *models.ProductDistributor) error {
	_, err := r.db.NewDelete().Model(product).Where("id = ?", product.ID).Exec(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductDistributorRepository) DeleteWithCondition(filter ProductDistributorRepositoryFilter) (int64, error) {
	var product models.ProductDistributor

	sl := r.db.NewDelete().Model(&product)
	sl = r.deleteQueryBuilder(sl, filter)

	res, err := sl.Exec(context.TODO())
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}
