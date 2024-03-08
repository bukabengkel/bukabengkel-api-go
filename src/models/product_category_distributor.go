package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ProductCategoryDistributor struct {
	bun.BaseModel `bun:"table:product_category_distributor,timestamps"` // Use bun.BaseModel for timestamps

	ID            uint64    `bun:"id,pk,autoIncrement"`
	ExternalID    string    `bun:"external_id"`
	DistributorID uint64    `bun:"distributor_id,fk:distributor_distributor_id"` // Corrected foreign key field
	Name          string    `bun:"name"`
	Code          string    `bun:"code"`
	Description   string    `bun:"description"`
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
	RemoteUpdate  bool      `bun:"remote_update"` // Use Go's boolean type

	Distributor *Distributor `bun:"rel:hasOne,to:distributor,foreignKey:distributor_id"` // Corrected model and table name
}
