package repository

import (
	"context"
	"fmt"

	"github.com/peang/bukabengkel-api-go/src/domain/entity"
	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/peang/bukabengkel-api-go/src/utils"
	"github.com/uptrace/bun"
)

type ProductRepository struct {
	db              *bun.DB
	imageRepository *ImageRepository
}

type ProductRepositoryFilter struct {
	Name    *string
	StoreID *int
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

	return query
}

func (r *ProductRepository) List(ctx context.Context, page int, perPage int, sort string, filter ProductRepositoryFilter) (*[]entity.Product, int, error) {
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
		return &[]entity.Product{}, count, nil
	}

	var entityProducts []entity.Product
	var entityProductIds []*int64
	for _, p := range products {
		entityProduct := models.LoadProductModel(p)

		entityProductIds = append(entityProductIds, p.ID)
		entityProducts = append(entityProducts, *entityProduct)
	}

	images, err := r.imageRepository.Find(ctx, 1, 5, "id", ImageRepositoryFilter{
		EntityIDS:  entityProductIds,
		EntityType: utils.Uint(1),
	})

	if err == nil {
		for key, p := range entityProducts {
			var productImages []*entity.Image
			for _, i := range images {
				if p.ID == i.EntityId {
					productImages = append(productImages, i)
				}
			}
			p.Images = productImages
			if len(productImages) > 0 {
				p.Thumbnail = productImages[0]
			}

			entityProducts[key] = p
		}
	}

	return &entityProducts, count, nil
}
