package models

import (
	"time"

	"github.com/peang/bukabengkel-api-go/src/domain/entity"
	"github.com/uptrace/bun"
)

type Location struct {
	bun.BaseModel `bun:"table:location"`

	ID        *int64    `bun:"id,pk"`
	Province  string    `bun:"province,notnull"`
	City      string    `bun:"city,notnull"`
	District  string    `bun:"district,notnull"`
	Urban     string    `bun:"urban,notnull"`
	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
}

func LoadLocationModel(l *Location) *entity.Location {
	if l != nil {
		return &entity.Location{
			ID:        *l.ID,
			Province:  l.Province,
			City:      l.City,
			District:  l.District,
			Urban:     l.Urban,
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
		}
	}

	return nil
}
