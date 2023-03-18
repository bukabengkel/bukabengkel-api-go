package entity

import "time"

type ImageType int

const (
	ImageProductType      ImageType = 1
	ImageProfileImageType ImageType = 2
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
	ID           int       `json:"id"`
	OwnerID      int       `json:"owner_id"`
	EntityId     int       `json:"entity_id"`
	EntityType   ImageType `json:"entity_type"`
	FileName     string    `json:"file_name"`
	OriginalName string    `json:"original_name"`
	Extension    string    `json:"extension"`
	Path         string    `json:"path"`
	Size         int       `json:"size"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
