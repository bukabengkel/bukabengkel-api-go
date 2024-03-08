package repository

import (
	"context"
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/uptrace/bun"
)

type ProductCategoryDistributorRepository struct {
	db *bun.DB
}

type ProductCategoryDistributorRepositoryFilter struct {
	Name *string
}

func NewProductCategoryDistributorRepository(db *bun.DB) *ProductCategoryDistributorRepository {
	return &ProductCategoryDistributorRepository{
		db: db,
	}
}

func (r *ProductCategoryDistributorRepository) queryBuilder(query *bun.SelectQuery, filter ProductCategoryDistributorRepositoryFilter) *bun.SelectQuery {
	if filter.Name != nil {
		query.Where("? ILIKE ?", bun.Ident("product_category_distributor.name"), fmt.Sprintf("%%%s%%", *filter.Name))
	}

	return query
}

func (r *ProductCategoryDistributorRepository) FindOne(filter ProductCategoryDistributorRepositoryFilter) (*models.ProductCategoryDistributor, error) {
	var productCategory models.ProductCategoryDistributor

	sl := r.db.NewSelect().Model(&productCategory)
	sl = r.queryBuilder(sl, filter)

	err := sl.Scan(context.TODO(), productCategory)
	if err != nil {
		return nil, err
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

func (r *ProductCategoryDistributorRepository) UpdateJob(productCategory *models.ProductCategoryDistributor) (*models.ProductCategoryDistributor, error) {
	_, err := r.db.NewUpdate().Model(productCategory).Where("id = ?", productCategory.ID).Exec(context.TODO())
	if err != nil {
		return nil, err
	}

	return productCategory, nil
}

func (r *ProductCategoryDistributorRepository) DeleteJob(productCategory *models.ProductCategoryDistributor) error {
	_, err := r.db.NewDelete().Model(productCategory).Where("id = ?", productCategory.ID).Exec(context.TODO())
	if err != nil {
		return err
	}

	return nil
}
