package entity

import "time"

type ImageEntityType int

const (
	ImageProductType      ImageEntityType = 1
	ImageProfileImageType ImageEntityType = 2
)

func (t ImageEntityType) String() string {
	switch t {
	case ImageProductType:
		return "product"
	case ImageProfileImageType:
		return "profile_image"
	default:
		return "unknown"
	}
}

// [{"id":4,"owner_id":1,"entity_id":1,"entity_type":2,"file_name":"c8bda027-5127-461a-9ba7-d8ff0985edac","original_name":"Screenshot 2023-01-20 at 16.45.42.png","extension":"image/png","path":"local/products/c8bda027-5127-461a-9ba7-d8ff0985eaac.png","size":179395,"created_at":"2023-02-05T09:13:02.556+00:00","updated_at":"2023-02-05T09:37:16.799+00:00"}]
type ImageEntity struct {
	ID           int             `json:"id"`
	OwnerID      int             `json:"owner_id"`
	EntityId     int             `json:"entity_id"`
	EntityType   ImageEntityType `json:"entity_type"`
	FileName     string          `json:"file_name"`
	OriginalName string          `json:"original_name"`
	Extension    string          `json:"extension"`
	Path         string          `json:"path"`
	Size         int             `json:"size"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
