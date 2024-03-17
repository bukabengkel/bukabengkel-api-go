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
}

type ProductRepositoryFilter struct {
	Name          *string
	StoreID       *int64
	StockMoreThan *uint
}

func NewProductRepository(db *bun.DB, imageRepository *ImageRepository) *ProductRepository {
	return &ProductRepository{
		db:              db,
		imageRepository: imageRepository,
	}
}

func (r *ProductRepository) queryBuilder(query *bun.SelectQuery, cond ProductRepositoryFilter) *bun.SelectQuery {
	if cond.Name != nil {
		query.Where("? ILIKE ?", bun.Ident("product.name"), fmt.Sprintf("%%%s%%", *cond.Name))
	}

	if cond.StoreID != nil {
		query.Where("? = ?", bun.Ident("product.store_id"), cond.StoreID)
	}

	if cond.StockMoreThan != nil {
		query.Where("? > ?", bun.Ident("product.stock"), cond.StockMoreThan)
	}

	return query
}

func (r *ProductRepository) List(ctx context.Context, page int, perPage int, sort string, filter ProductRepositoryFilter) (*[]models.Product, int, error) {
	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimit(page, perPage)

	var products []models.Product
	sl := r.db.NewSelect().Model(&products)
	sl = r.queryBuilder(sl, filter)

	count, err := sl.
		Relation("Store").
		Relation("Brand").
		Relation("Category").
		Limit(limit).Offset(offset).OrderExpr(sorts).ScanAndCount(context.TODO())
	if err != nil {
		return nil, 0, err
	}

	if len(products) == 0 {
		return &[]models.Product{}, count, nil
	}

	var entityProducts []models.Product
	var entityProductIds []*uint64
	for _, p := range products {
		entityProductIds = append(entityProductIds, p.ID)
		entityProducts = append(entityProducts, p)
	}

	images, err := r.imageRepository.Find(ctx, 1, 50, "id", ImageRepositoryFilter{
		EntityIDS:  entityProductIds,
		EntityType: utils.Uint(1),
	})

	if err == nil {
		for key, p := range entityProducts {
			var productImages []models.Image
			for _, img := range images {
				if img.EntityID == *p.ID {
					img.Path = r.imageRepository.s3service.BuildUrl(img.Path, 500, 500)
					productImages = append(productImages, img)
				}
			}

			p.Images = productImages
			entityProducts[key] = p
		}
	}

	return &entityProducts, count, nil
}
