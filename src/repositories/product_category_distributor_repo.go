package repository

import (
	"context"
	"database/sql"

	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/uptrace/bun"
)

type ProductCategoryDistributorRepository struct {
	db *bun.DB
}

type ProductCategoryDistributorRepositoryFilter struct {
	ExternalID    *string
	DistributorID *uint64
	Code          *string
	RemoteUpdate  *bool
}

type ProductCategoryDistributorRepositoryValues struct {
	RemoteUpdate *bool
}

func NewProductCategoryDistributorRepository(db *bun.DB) *ProductCategoryDistributorRepository {
	return &ProductCategoryDistributorRepository{
		db: db,
	}
}

func (r *ProductCategoryDistributorRepository) queryBuilder(
	query *bun.SelectQuery,
	filter ProductCategoryDistributorRepositoryFilter,
) *bun.SelectQuery {
	if filter.DistributorID != nil {
		query.Where("? = ?", bun.Ident("distributor_id"), filter.DistributorID)
	}

	if filter.ExternalID != nil {
		query.Where("? = ?", bun.Ident("external_id"), filter.ExternalID)
	}

	if filter.Code != nil {
		query.Where("? = ?", bun.Ident("code"), filter.Code)
	}

	return query
}

func (r *ProductCategoryDistributorRepository) updateQueryBuilder(
	query *bun.UpdateQuery,
	filter ProductCategoryDistributorRepositoryFilter,
	values ProductCategoryDistributorRepositoryValues,
) *bun.UpdateQuery {
	if filter.DistributorID != nil {
		query.Where("? = ?", bun.Ident("distributor_id"), filter.DistributorID)
	}

	if filter.ExternalID != nil {
		query.Where("? = ?", bun.Ident("external_id"), filter.ExternalID)
	}

	if filter.Code != nil {
		query.Where("? = ?", bun.Ident("code"), filter.Code)
	}

	if values.RemoteUpdate != nil {
		query.Set("? = ?", bun.Ident("remote_update"), values.RemoteUpdate)
	}

	return query
}

func (r *ProductCategoryDistributorRepository) deleteQueryBuilder(
	query *bun.DeleteQuery,
	filter ProductCategoryDistributorRepositoryFilter,
) *bun.DeleteQuery {
	if filter.DistributorID != nil {
		query.Where("? = ?", bun.Ident("distributor_id"), filter.DistributorID)
	}

	if filter.ExternalID != nil {
		query.Where("? = ?", bun.Ident("external_id"), filter.ExternalID)
	}

	if filter.Code != nil {
		query.Where("? = ?", bun.Ident("code"), filter.Code)
	}

	if filter.RemoteUpdate != nil {
		query.Where("? = ?", bun.Ident("remote_update"), filter.RemoteUpdate)
	}

	return query
}

func (r *ProductCategoryDistributorRepository) FindOne(filter ProductCategoryDistributorRepositoryFilter) (*models.ProductCategoryDistributor, error) {
	var productCategory models.ProductCategoryDistributor

	sl := r.db.NewSelect().Model(&productCategory)
	sl = r.queryBuilder(sl, filter)

	err := sl.Scan(context.TODO())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &productCategory, nil
}

func (r *ProductCategoryDistributorRepository) Save(productCategory *models.ProductCategoryDistributor) (*models.ProductCategoryDistributor, error) {
	_, err := r.db.NewInsert().Model(productCategory).Returning("id").Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	err = r.db.NewSelect().Model(productCategory).WherePK().Scan(context.TODO())
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

func (r *ProductCategoryDistributorRepository) Update(productCategory *models.ProductCategoryDistributor) (*models.ProductCategoryDistributor, error) {
	_, err := r.db.NewUpdate().Model(productCategory).Where("id = ?", productCategory.ID).Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

func (r *ProductCategoryDistributorRepository) UpdateWithCondition(filter ProductCategoryDistributorRepositoryFilter, values ProductCategoryDistributorRepositoryValues) (int64, error) {
	var productCategory models.ProductCategoryDistributor

	sl := r.db.NewUpdate().Model(&productCategory)
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

func (r *ProductCategoryDistributorRepository) Delete(productCategory *models.ProductCategoryDistributor) error {
	_, err := r.db.NewDelete().Model(productCategory).Where("id = ?", productCategory.ID).Exec(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductCategoryDistributorRepository) DeleteWithCondition(filter ProductCategoryDistributorRepositoryFilter) (int64, error) {
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
