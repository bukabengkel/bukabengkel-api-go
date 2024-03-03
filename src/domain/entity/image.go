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
	ID           int64
	EntityId     int64
	EntityType   ImageType
	FileName     string
	OriginalName string
	Extension    string
	Path         string
	Size         int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
