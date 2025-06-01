package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Geolocation struct {
	Lat  float64
	Long float64
}

type Store struct {
	bun.BaseModel `bun:"table:store"`

	ID             *uint64     `bun:"id,pk"`
	Key            string      `bun:"key,notnull,unique"`
	Name           string      `bun:"name,notnull"`
	Type           uint8       `bun:"type,notnull"`
	LocationDetail string      `bun:"location_detail"`
	Geolocation    Geolocation `bun:"geolocation"`
	CreatedAt      time.Time   `bun:"created_at"`
	CreatedBy      string      `bun:"created_by"`
	UpdatedAt      time.Time   `bun:"updated_at"`
	UpdatedBy      string      `bun:"updated_by"`
	DeletedAt      time.Time   `bun:"deleted_at"`
	DeletedBy      string      `bun:"deleted_by"`
}
