package entity

import "time"

type ImageEntityType int

const (
	Product      ImageEntityType = 1
	ProfileImage ImageEntityType = 2
)

func (t ImageEntityType) String() string {
	switch t {
	case Product:
		return "product"
	case ProfileImage:
		return "profile_image"
	default:
		return "unknown"
	}
}

type ImageEntity struct {
	ID           int
	ImageEntity  ImageEntityType
	OwnerID      int
	FileName     string
	OriginalName string
	Extension    string
	Path         string
	Size         int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
