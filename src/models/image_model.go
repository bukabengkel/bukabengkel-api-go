package models

import (
	"time"

	"github.com/peang/bukabengkel-api-go/src/domain/entity"
	"github.com/uptrace/bun"
)

type Image struct {
	bun.BaseModel `bun:"table:image"`

	ID           *int64    `bun:"id,pk"`
	OwnerID      int64     `bun:"owner_id,notnull"`
	EntityID     int64     `bun:"entity_id,notnull"`
	EntityType   int64     `bun:"entity_type,notnull"`
	FileName     string    `bun:"file_name,notnull"`
	OriginalName string    `bun:"original_name,notnull"`
	Extension    string    `bun:"extension,notnull"`
	Path         string    `bun:"path,notnull"`
	Size         int       `bun:"size,notnull"`
	CreatedAt    time.Time `bun:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at"`
}

func LoadImageModel(i *Image) *entity.Image {
	if i != nil {
		return &entity.Image{
			ID:           *i.ID,
			EntityId:     i.EntityID,
			EntityType:   entity.ImageType(i.EntityType),
			FileName:     i.FileName,
			OriginalName: i.OriginalName,
			Extension:    i.Extension,
			Path:         i.Path,
			Size:         i.Size,
			CreatedAt:    i.CreatedAt,
			UpdatedAt:    i.UpdatedAt,
		}
	}

	return nil
}

func LoadImageModels(i []Image) *[]entity.Image {
	var entityImages []entity.Image

	for _, i := range i {
		entityImage := entity.Image{
			ID:           *i.ID,
			EntityId:     i.EntityID,
			EntityType:   entity.ImageType(i.EntityType),
			FileName:     i.FileName,
			OriginalName: i.OriginalName,
			Extension:    i.Extension,
			Path:         i.Path,
			Size:         i.Size,
			CreatedAt:    i.CreatedAt,
			UpdatedAt:    i.UpdatedAt,
		}

		entityImages = append(entityImages, entityImage)
	}

	return &entityImages
}
