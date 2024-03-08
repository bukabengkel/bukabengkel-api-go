package repository

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/models"
	file_service "github.com/peang/bukabengkel-api-go/src/services/file_services"
	"github.com/peang/bukabengkel-api-go/src/utils"
	"github.com/uptrace/bun"
)

type ImageRepository struct {
	db          *bun.DB
	fileService *file_service.S3Service
}

type ImageRepositoryFilter struct {
	EntityID   *uint64
	EntityIDS  []*uint64
	EntityType *uint
}

func NewImageRepository(db *bun.DB, fileService *file_service.S3Service) *ImageRepository {
	return &ImageRepository{
		db:          db,
		fileService: fileService,
	}
}

func (r *ImageRepository) queryBuilder(query *bun.SelectQuery, cond ImageRepositoryFilter) *bun.SelectQuery {
	if cond.EntityID != nil {
		query.Where("? = ?", bun.Ident("entity_id"), cond.EntityID)
	}

	if cond.EntityIDS != nil {
		query.Where("? IN (?)", bun.Ident("entity_id"), bun.In(cond.EntityIDS))
	}

	if cond.EntityType != nil {
		query.Where("? = ?", bun.Ident("entity_type"), cond.EntityType)
	}

	return query
}

func (r *ImageRepository) Find(ctx context.Context, page int, perPage int, sort string, filter ImageRepositoryFilter) ([]*models.Image, error) {
	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimit(page, perPage)

	var images []models.Image
	sl := r.db.NewSelect().Model(&images)
	sl = r.queryBuilder(sl, filter)

	err := sl.Limit(limit).Offset(offset).OrderExpr(sorts).Scan(ctx)
	if err != nil {
		return nil, err
	}

	if len(images) == 0 {
		return []*models.Image{}, nil
	}

	var entityImages []*models.Image

	for _, image := range images {
		// entityImage := models.LoadImageModel(&image)
		image.Path = r.fileService.BuildUrl(image.Path, 500, 500)
		entityImages = append(entityImages, &image)
	}

	return entityImages, nil
}

func (r *ImageRepository) FindAndCount(ctx context.Context, page int, perPage int, sort string, filter ImageRepositoryFilter) (*[]models.Image, int, error) {
	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimit(page, perPage)

	var images []models.Image
	sl := r.db.NewSelect().Model(&images)
	sl = r.queryBuilder(sl, filter)

	count, err := sl.
		Limit(limit).Offset(offset).OrderExpr(sorts).ScanAndCount(context.TODO())
	if err != nil {
		return nil, 0, err
	}

	if len(images) == 0 {
		return &[]models.Image{}, count, nil
	}

	var entityImages []models.Image

	for _, image := range images {
		// entityImage := models.LoadImageModel(&image)

		entityImages = append(entityImages, image)
	}

	return &entityImages, count, nil
}
