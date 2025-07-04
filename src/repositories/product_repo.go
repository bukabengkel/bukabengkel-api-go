package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/peang/bukabengkel-api-go/src/services/cache_services"
	"github.com/peang/bukabengkel-api-go/src/utils"
	"github.com/uptrace/bun"
)

type ProductRepository struct {
	db              *bun.DB
	cacheService    cache_services.CacheService
	imageRepository *ImageRepository
	statements      *productStatements
}

type ProductRepositoryFilter struct {
	Name          *string
	CategoryId    *string
	StoreID       *uint
	StockMoreThan *uint
	Status        *uint
	StartDate     *time.Time
	EndDate       *time.Time
}

// Separate structure to hold prepared statements
type productStatements struct {
	// Base queries
	listQuery    *bun.SelectQuery
	findOneQuery *bun.SelectQuery

	// Common conditions
	nameCondition     *bun.SelectQuery
	categoryCondition *bun.SelectQuery
	storeCondition    *bun.SelectQuery
	stockCondition    *bun.SelectQuery
	statusCondition   *bun.SelectQuery
}

// Initialize prepared statements with components
func newProductStatements(db *bun.DB) *productStatements {
	baseModel := (*models.Product)(nil)
	relations := []string{"Store", "Brand", "Category"}

	// Prepare base queries
	listQuery := db.NewSelect().Model(baseModel)
	for _, relation := range relations {
		listQuery = listQuery.Relation(relation)
	}

	// Prepare common conditions
	nameCondition := db.NewSelect().Where("? ILIKE ?", bun.Ident("product.name"), "")
	categoryCondition := db.NewSelect().Where("? = ?", bun.Ident("product.category_id"), "")
	storeCondition := db.NewSelect().Where("? = ?", bun.Ident("product.store_id"), 0)
	stockCondition := db.NewSelect().Where("? > ?", bun.Ident("product.stock"), 0)
	statusCondition := db.NewSelect().Where("? = ?", bun.Ident("product.status"), 0)

	return &productStatements{
		listQuery:         listQuery,
		findOneQuery:      db.NewSelect().Model(baseModel),
		nameCondition:     nameCondition,
		categoryCondition: categoryCondition,
		storeCondition:    storeCondition,
		stockCondition:    stockCondition,
		statusCondition:   statusCondition,
	}
}

func NewProductRepository(db *bun.DB, imageRepository *ImageRepository, cacheService cache_services.CacheService) *ProductRepository {
	return &ProductRepository{
		db:              db,
		imageRepository: imageRepository,
		statements:      newProductStatements(db),
		cacheService:    cacheService,
	}
}

func (r *ProductRepository) queryBuilder(query *bun.SelectQuery, cond ProductRepositoryFilter) *bun.SelectQuery {
	if cond.Name != nil {
		query = query.Where("? ILIKE ?", bun.Ident("product.name"), fmt.Sprintf("%%%s%%", *cond.Name))
	}

	if cond.CategoryId != nil {
		query = query.Where("? = ?", bun.Ident("product.category_id"), cond.CategoryId)
	}

	if cond.StoreID != nil {
		query = query.Where("? = ?", bun.Ident("product.store_id"), *cond.StoreID)
	}

	if cond.StockMoreThan != nil {
		query = query.Where("? > ?", bun.Ident("product.stock"), cond.StockMoreThan)
	}

	if cond.Status != nil {
		query = query.Where("? = ?", bun.Ident("product.status"), cond.Status)
	}

	if cond.StartDate != nil && cond.EndDate != nil {
		query = query.Where("? BETWEEN ? AND ?", bun.Ident("product.created_at"), cond.StartDate, cond.EndDate)
	}

	return query
}

func (r *ProductRepository) List(ctx context.Context, page int, perPage int, sort string, filter ProductRepositoryFilter) (*[]models.Product, int, error) {
  // // Generate cache key berdasarkan parameter
	// cacheKey := r.generateCacheKey("product_list", page, perPage, sort, filter)

	// // Coba ambil dari cache terlebih dahulu
	// if cached, err := r.cacheService.Get(ctx, cacheKey); err == nil {
	// 	if result, ok := cached.(*[]models.Product); ok {
	// 		return result, 0, nil
	// 	}
	// }

	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimit(page, perPage)

	var products []models.Product

	// Use prepared statement
	query := r.statements.listQuery
	query = r.queryBuilder(query, filter)

	count, err := query.
		Limit(limit).
		Offset(offset).
		OrderExpr(sorts).
		ScanAndCount(ctx, &products)

	if err != nil {
		return nil, 0, err
	}

	if len(products) == 0 {
		return &[]models.Product{}, count, nil
	}

	// Optimize image loading with batch processing
	if err := r.loadProductImages(ctx, &products); err != nil {
		return nil, 0, err
	}

	// // Simpan ke cache dengan TTL yang sesuai
	// err = r.cacheService.Set(ctx, cacheKey, &products, time.Second*30)
	// if err != nil {
	// 	fmt.Println("Error setting cache", err)
	// }

	return &products, count, nil
}

func (r *ProductRepository) ListV2(ctx context.Context, page int, perPage int, sort string, filter ProductRepositoryFilter) (*[]models.Product, bool, error) {
	// // Generate cache key berdasarkan parameter
	// cacheKey := r.generateCacheKey("product_list_v2", page, perPage, sort, filter)

	// // Coba ambil dari cache terlebih dahulu
	// if cached, err := r.cacheService.Get(ctx, cacheKey); err == nil {
	// 	if result, ok := cached.(*[]models.Product); ok {
	// 		return result, false, nil
	// 	}
	// }

	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimitV2(page, perPage)

	var products []models.Product
	var hasNext bool = false

	// Use prepared statement
	query := r.statements.listQuery
	query = r.queryBuilder(query, filter)

	err := query.
		Limit(limit).
		Offset(offset).
		OrderExpr(sorts).
		Scan(ctx, &products)

	if err != nil {
		return nil, false, err
	}

	if len(products) == 0 {
		return &[]models.Product{}, false, nil
	}

	// Optimize image loading with batch processing
	if err := r.loadProductImages(ctx, &products); err != nil {
		return nil, false, err
	}

	if len(products) > perPage {
		hasNext = true
		products = products[:len(products)-1]
	}

	// Simpan ke cache dengan TTL yang sesuai
	// err = r.cacheService.Set(ctx, cacheKey, &products, time.Second*30)
	// if err != nil {
	// 	fmt.Println("Error setting cache", err)
	// }

	return &products, hasNext, nil
}

func (r *ProductRepository) Count(ctx context.Context, filter ProductRepositoryFilter) (int, error) {
	sl := r.db.NewSelect().Table("product")
	sl = r.queryBuilder(sl, filter)

	count, err := sl.Count(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Helper method to load product images efficiently
func (r *ProductRepository) loadProductImages(ctx context.Context, products *[]models.Product) error {
	if len(*products) == 0 {
		return nil
	}

	entityProductIds := make([]uint64, 0, len(*products))
	for _, p := range *products {
		entityProductIds = append(entityProductIds, p.ID)
	}

	images, err := r.imageRepository.Find(ctx, 1, 50, "id", ImageRepositoryFilter{
		EntityIDS:  &entityProductIds,
		EntityType: utils.Uint64(1),
	})

	if err != nil {
		return err
	}

	// Create image map for efficient lookup
	imageMap := make(map[uint64][]models.Image)
	for _, img := range images {
		img.Path = r.imageRepository.fileService.BuildUrl(img.Path, 500, 500)
		imageMap[img.EntityID] = append(imageMap[img.EntityID], img)
	}

	// Assign images to products
	for i := range *products {
		if imgs, ok := imageMap[(*products)[i].ID]; ok {
			(*products)[i].Images = imgs
		}
	}

	return nil
}

func (r *ProductRepository) generateCacheKey(prefix string, page int, perPage int, sort string, filter ProductRepositoryFilter) string {
	key := fmt.Sprintf("%s:%d:%d:%s", prefix, page, perPage, sort)
	if filter.StoreID != nil {
		key += fmt.Sprintf(":store_%d", *filter.StoreID)
	}
	if filter.Name != nil {
		key += fmt.Sprintf(":name_%s", *filter.Name)
	}
	if filter.CategoryId != nil {
		key += fmt.Sprintf(":cat_%s", *filter.CategoryId)
	}
	if filter.StockMoreThan != nil {
		key += fmt.Sprintf(":stock_%d", *filter.StockMoreThan)
	}
	if filter.Status != nil {
		key += fmt.Sprintf(":status_%d", *filter.Status)
	}

	return key
}
