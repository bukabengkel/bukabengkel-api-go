package repository

import (
	"context"
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/peang/bukabengkel-api-go/src/utils"
	"github.com/uptrace/bun"
)

type ProductRepository struct {
	db              *bun.DB
	imageRepository *ImageRepository
	statements      *productStatements
}

type ProductRepositoryFilter struct {
	Name          *string
	CategoryId    *string
	StoreID       *uint
	StockMoreThan *uint
	Status        *uint
}

// Separate structure to hold prepared statements
type productStatements struct {
	// Base model and relations for list query
	listModel     interface{}
	listRelations []string
}

// Initialize prepared statements with components rather than full queries
func newProductStatements(db *bun.DB) *productStatements {
	return &productStatements{
		listModel:     (*models.Product)(nil),
		listRelations: []string{"Store", "Brand", "Category"},
	}
}

func NewProductRepository(db *bun.DB, imageRepository *ImageRepository) *ProductRepository {
	return &ProductRepository{
		db:              db,
		imageRepository: imageRepository,
		statements:      newProductStatements(db),
	}
}

func (r *ProductRepository) queryBuilder(query *bun.SelectQuery, cond ProductRepositoryFilter) *bun.SelectQuery {
	if cond.Name != nil {
		query.Where("? ILIKE ?", bun.Ident("product.name"), fmt.Sprintf("%%%s%%", *cond.Name))
	}

	if cond.CategoryId != nil {
		query.Where("? = ?", bun.Ident("product.category_id"), cond.CategoryId)
	}

	if cond.StoreID != nil {
		query.Where("? = ?", bun.Ident("product.store_id"), *cond.StoreID)
	}

	if cond.StockMoreThan != nil {
		query.Where("? > ?", bun.Ident("product.stock"), cond.StockMoreThan)
	}

	if cond.Status != nil {
		query.Where("? = ?", bun.Ident("product.status"), cond.Status)
	}

	return query
}

func (r *ProductRepository) List(ctx context.Context, page int, perPage int, sort string, filter ProductRepositoryFilter) (*[]models.Product, int, error) {
	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimit(page, perPage)

	var products []models.Product

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
		ScanAndCount(ctx)

	if err != nil {
		return nil, 0, err
	}

	if len(products) == 0 {
		return &[]models.Product{}, count, nil
	}

	entityProductIds := make([]uint64, 0, len(products))
	for _, p := range products {
		entityProductIds = append(entityProductIds, p.ID)
	}

	images, err := r.imageRepository.Find(ctx, 1, 50, "id", ImageRepositoryFilter{
		EntityIDS:  &entityProductIds,
		EntityType: utils.Uint64(1),
	})

	if err == nil {
		imageMap := make(map[uint64][]models.Image)
		for _, img := range images {
			img.Path = r.imageRepository.fileService.BuildUrl(img.Path, 500, 500)
			imageMap[img.EntityID] = append(imageMap[img.EntityID], img)
		}

		for i := range products {
			if imgs, ok := imageMap[products[i].ID]; ok {
				products[i].Images = imgs
			}
		}
	}

	return &products, count, nil
}

func (r *ProductRepository) ListV2(ctx context.Context, page int, perPage int, sort string, filter ProductRepositoryFilter) (*[]models.Product, bool, error) {
	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimitV2(page, perPage)

	var products []models.Product
	var hasNext bool = false

	// Create a new query using the base components
	sl := r.db.NewSelect().Model(&products)

	// Add relations
	for _, relation := range r.statements.listRelations {
		sl = sl.Relation(relation)
	}

	sl = r.queryBuilder(sl, filter)

	err := sl.
		Limit(limit).
		Offset(offset).
		OrderExpr(sorts).
		Scan(ctx)

	if err != nil {
		return nil, false, err
	}

	if len(products) == 0 {
		return &[]models.Product{}, false, nil
	}

	entityProductIds := make([]uint64, 0, len(products))
	for _, p := range products {
		entityProductIds = append(entityProductIds, p.ID)
	}

	images, err := r.imageRepository.Find(ctx, 1, 50, "id", ImageRepositoryFilter{
		EntityIDS:  &entityProductIds,
		EntityType: utils.Uint64(1),
	})

	if err == nil {
		imageMap := make(map[uint64][]models.Image)
		for _, img := range images {
			img.Path = r.imageRepository.fileService.BuildUrl(img.Path, 500, 500)
			imageMap[img.EntityID] = append(imageMap[img.EntityID], img)
		}

		for i := range products {
			if imgs, ok := imageMap[products[i].ID]; ok {
				products[i].Images = imgs
			}
		}
	}

	if len(products) > perPage {
		hasNext = true
		products = products[:len(products)-1]
	}

	return &products, hasNext, nil
}
