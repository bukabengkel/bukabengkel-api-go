package entity

import "time"

type Store struct {
	ID             int
	Key            string
	Name           string
	Type           int
	Location       Location
	LocationDetail string
	Geolocation    struct {
		Lat  float64
		Long float64
	}
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt time.Time
	DeletedBy string
}
