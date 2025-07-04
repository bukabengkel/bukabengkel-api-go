package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ImageType int

const (
	ImageProductType      ImageType = 1
	ImageProfileImageType ImageType = 2
)

const (
	ImageProduct            = "products"
	ImageProductDistributor = "products_distributor"
)

func (t ImageType) String() string {
	switch t {
	case ImageProductType:
		return "product"
	case ImageProfileImageType:
		return "profile_image"
	default:
		return "unknown"
	}
}

type Image struct {
	bun.BaseModel `bun:"table:image"`

	ID           uint64    `bun:"id,pk,nullzero"`
	OwnerID      uint64    `bun:"owner_id,notnull"`
	EntityID     uint64    `bun:"entity_id,notnull"`
	EntityType   uint64    `bun:"entity_type,notnull"`
	FileName     string    `bun:"file_name,notnull"`
	OriginalName string    `bun:"original_name,notnull"`
	Extension    string    `bun:"extension,notnull"`
	Path         string    `bun:"path,notnull"`
	Size         int       `bun:"size,notnull"`
	CreatedAt    time.Time `bun:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at"`
}
