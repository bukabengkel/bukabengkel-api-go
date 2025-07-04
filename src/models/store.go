package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Store struct {
	bun.BaseModel `bun:"table:store"`

	ID             *uint64   `bun:"id,pk"`
	Key            string    `bun:"key,notnull,unique"`
	Name           string    `bun:"name,notnull"`
	Type           uint8     `bun:"type,notnull"`
	LocationDetail string    `bun:"location_detail"`
	Geolocation    string    `bun:"geolocation,type:geometry(POINT,3857)"`
	CreatedAt      time.Time `bun:"created_at"`
	CreatedBy      string    `bun:"created_by"`
	UpdatedAt      time.Time `bun:"updated_at"`
	UpdatedBy      string    `bun:"updated_by"`
	DeletedAt      time.Time `bun:"deleted_at"`
	DeletedBy      string    `bun:"deleted_by"`
}
