package models

import (
	"time"

	"github.com/peang/bukabengkel-api-go/src/domain/entity"
	"github.com/uptrace/bun"
)

type Geolocation struct {
	Lat  float64
	Long float64
}

type Store struct {
	bun.BaseModel `bun:"table:store"`

	ID             *int64      `bun:"id,pk"`
	Key            string      `bun:"key,notnull,unique"`
	Name           string      `bun:"name,notnull"`
	Type           int         `bun:"type,notnull"`
	LocationID     *int64      `bun:"location_id,notnull"`
	Location       *Location   `bun:"rel:belongs-to,join:location_id=id"`
	LocationDetail string      `bun:"location_detail"`
	Geolocation    Geolocation `bun:"geolocation"`
	CreatedAt      time.Time   `bun:"created_at"`
	CreatedBy      string      `bun:"created_by"`
	UpdatedAt      time.Time   `bun:"updated_at"`
	UpdatedBy      string      `bun:"updated_by"`
	DeletedAt      time.Time   `bun:"deleted_at"`
	DeletedBy      string      `bun:"deleted_by"`
}

func LoadStoreModel(s *Store) *entity.Store {
	if s != nil {
		return &entity.Store{
			ID:             *s.ID,
			Key:            s.Key,
			Name:           s.Name,
			Type:           s.Type,
			Location:       LoadLocationModel(s.Location),
			LocationDetail: s.LocationDetail,
			Geolocation:    s.Geolocation,
			CreatedAt:      s.CreatedAt,
			CreatedBy:      s.CreatedBy,
			UpdatedAt:      s.UpdatedAt,
			UpdatedBy:      s.UpdatedBy,
			DeletedAt:      s.DeletedAt,
			DeletedBy:      s.DeletedBy,
		}
	}

	return nil
}
