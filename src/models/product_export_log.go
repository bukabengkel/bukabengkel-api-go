package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ProductExportStatus int

const (
	LogDraft ProductExportStatus = iota
	LogOnProgress
	LogDone
	LogError
)

func (s ProductExportStatus) String() string {
	switch s {
	case LogDraft:
		return "draft"
	case LogOnProgress:
		return "on-progress"
	case LogDone:
		return "done"
	case LogError:
		return "error"
	default:
		return "unknown"
	}
}

type ProductExportLog struct {
	bun.BaseModel `bun:"table:product_export_log"`

	ID        uint64              `bun:"id,pk,nullzero"`
	StoreID   uint64              `bun:"store_id,notnull"`
	UserID    uint64              `bun:"user_id,notnull"`
	Status    ProductExportStatus `bun:"status,notnull"`
	PathFile  string              `bun:"path_file,type:text"`
	DoneAt    time.Time           `bun:"done_at,nullzero"`
	CreatedAt time.Time           `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time           `bun:"updated_at"`
}
