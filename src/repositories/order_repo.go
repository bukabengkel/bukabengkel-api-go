package repository

import (
	"context"
	"fmt"
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
		query.Where("? BETWEEN ? AND ?", bun.Ident("order_date"), cond.StartDate, cond.EndDate)
	}

	return query
}

func (r *OrderRepository) CountReport(ctx context.Context, filter OrderRepositoryFilter) (*struct {
	totalSales float32
}, error) {
	var report struct {
		TotalSales float32 `bun:"total"`
	}
	sl := r.db.NewSelect().Table("order")
	sl = r.queryBuilder(sl, filter)

	err := sl.ColumnExpr("SUM(total) as total_sales").Scan(ctx, &report)
	if err != nil {
		return nil, err
	}

	fmt.Println(report)
	return &struct{ totalSales float32 }{
		totalSales: report.TotalSales,
	}, nil
}
