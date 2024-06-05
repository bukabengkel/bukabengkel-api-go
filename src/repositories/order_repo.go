package repository

import (
	"time"

	"github.com/uptrace/bun"
)

type OrderRepository struct {
	db *bun.DB
}

type OrderRepositoryFilter struct {
	StoreID   *uint
	StartDate *time.Time
	EndDate   *time.Time
}

func NewOrderRepository(db *bun.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) queryBuilder(query *bun.SelectQuery, cond OrderRepositoryFilter) *bun.SelectQuery {
	if cond.StoreID != nil {
		query.Where("? = ?", bun.Ident("store_id"), *cond.StoreID)
	}

	if cond.StartDate != nil && cond.EndDate != nil {
		query.Where("? BETWEEN ? AND ?", bun.Ident("created_at"), cond.StartDate, cond.EndDate)
	}

	return query
}
