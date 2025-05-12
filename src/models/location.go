package models

import (
	"time"

	"github.com/uptrace/bun"
)

type RajaOngkirLocation struct {
	bun.BaseModel `bun:"table:location_raja_ongkir"`

	ID              uint64    `bun:"id,pk"`
	EntityType      string    `bun:"entity_type"`
	EntityID        uint64    `bun:"entity_id"`
	LocationID      uint64    `bun:"location_id"`
	ProvinceName    string    `bun:"province_name"`
	CityName        string    `bun:"city_name"`
	DistrictName    string    `bun:"district_name"`
	SubdistrictName string    `bun:"subdistrict_name"`
	ZipCode         string    `bun:"zip_code"`
	Label           string    `bun:"label"`
	CreatedAt       time.Time `bun:"created_at"`
	UpdatedAt       time.Time `bun:"updated_at"`
}