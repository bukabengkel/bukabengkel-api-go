package repository

import (
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

// func (r *ProductCategoryDistributorRepository) queryBuilder(query *bun.SelectQuery, cond ProductCategoryDistributorRepositoryFilter) *bun.SelectQuery {
// 	if cond.Name != nil {
// 		query.Where("? ILIKE ?", bun.Ident("product_category_distributor.name"), fmt.Sprintf("%%%s%%", *cond.Name))
// 	}

// 	return query
// }
