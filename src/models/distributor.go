package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Distributor struct {
	bun.BaseModel `bun:"table:distributor,timestamps"` // Use bun.BaseModel for timestamps

	ID             uint64    `bun:"id,pk,autoIncrement"`
	Key            string    `bun:"key"`
	Name           string    `bun:"name"`
	LocationID     uint64    `bun:"location_id,fk:location_location_id"` // Corrected foreign key field
	LocationDetail string    `bun:"location_detail"`
	CreatedAt      time.Time `bun:"created_at"`
	UpdatedAt      time.Time `bun:"updated_at"`

	Location *Location `bun:"rel:hasOne,to:location,foreignKey:location_id"` // Corrected model and table name
}
