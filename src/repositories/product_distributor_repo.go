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
	statements  *productDistributorStatements
}

// Separate structure to hold prepared statements
type productDistributorStatements struct {
	// Base models and relations for queries
	listModel     interface{}
	listRelations []string
	findOneModel  interface{}
	findOneRelations []string
}

// Initialize prepared statements with components
func newProductDistributorStatements(db *bun.DB) *productDistributorStatements {
	return &productDistributorStatements{
		listModel:     (*models.ProductDistributor)(nil),
		listRelations: []string{"Distributor", "Category"},
		findOneModel:  (*models.ProductDistributor)(nil),
		findOneRelations: []string{"Distributor", "Category"},
	}
}

type ProductDistributorRepositoryFilter struct {
	ID            *string
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
		statements:  newProductDistributorStatements(db),
	}
}

func (r *ProductDistributorRepository) queryBuilder(query *bun.SelectQuery, filter ProductDistributorRepositoryFilter) *bun.SelectQuery {
	if filter.ID != nil {
		query.Where("? = ?", bun.Ident("product_distributor.key"), filter.ID)
	}

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
	
	// Create a new query using the base components
	sl := r.db.NewSelect().Model(&products)
	
	// Add relations
	for _, relation := range r.statements.listRelations {
		sl = sl.Relation(relation)
	}
	
	sl = r.queryBuilder(sl, filter)

	count, err := sl.
		Limit(limit).
		Offset(offset).
		OrderExpr(sorts).
		ScanAndCount(ctx) // Use the passed context
		
	if err != nil {
		return nil, 0, err
	}

	if len(products) == 0 {
		return &[]models.ProductDistributor{}, count, nil
	}

	// Pre-allocate the result slice with the right capacity
	entityProducts := make([]models.ProductDistributor, 0, len(products))
	for _, p := range products {
		p.ThumbnailCDN = r.fileService.BuildUrl(p.Thumbnail, 200, 200)
		entityProducts = append(entityProducts, p)
	}

	return &entityProducts, count, nil
}

func (r *ProductDistributorRepository) Save(ctx context.Context, product *models.ProductDistributor) (*models.ProductDistributor, error) {
	_, err := r.db.NewInsert().Model(product).Returning("id").Exec(ctx)
	if err != nil {
		return nil, err
	}

	err = r.db.NewSelect().Model(product).WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductDistributorRepository) Update(ctx context.Context, product *models.ProductDistributor) (*models.ProductDistributor, error) {
	_, err := r.db.NewUpdate().Model(product).Where("id = ?", product.ID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductDistributorRepository) UpdateWithCondition(ctx context.Context, filter ProductDistributorRepositoryFilter, values ProductDistributorRepositoryValues) (int64, error) {
	var product models.ProductDistributor

	sl := r.db.NewUpdate().Model(&product)
	sl = r.updateQueryBuilder(sl, filter, values)

	res, err := sl.Exec(ctx)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ProductDistributorRepository) FindOne(ctx context.Context, filter ProductDistributorRepositoryFilter) (*models.ProductDistributor, error) {
	var product models.ProductDistributor

	// Create a new query using the base components
	sl := r.db.NewSelect().Model(&product)
	
	// Add relations
	for _, relation := range r.statements.findOneRelations {
		sl = sl.Relation(relation)
	}
	
	sl = r.queryBuilder(sl, filter)

	err := sl.Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	product.ThumbnailCDN = r.fileService.BuildUrl(product.Thumbnail, 200, 200)
	return &product, nil
}

func (r *ProductDistributorRepository) Delete(ctx context.Context, product *models.ProductDistributor) error {
	_, err := r.db.NewDelete().Model(product).Where("id = ?", product.ID).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductDistributorRepository) DeleteWithCondition(ctx context.Context, filter ProductDistributorRepositoryFilter) (int64, error) {
	var product models.ProductDistributor

	sl := r.db.NewDelete().Model(&product)
	sl = r.deleteQueryBuilder(sl, filter)

	res, err := sl.Exec(ctx)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}
