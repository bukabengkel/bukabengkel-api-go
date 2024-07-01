package repository

import (
	"context"
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
	TotalSales   float32
	TotalNett    float32
	TotalProduct int
}, error) {
	var report struct {
		TotalSales   float32
		TotalNett    float32
		TotalProduct int
	}
	sl := r.db.NewSelect().Table("order").Join(`LEFT JOIN order_item as oi ON "oi".id = "order".id`)
	sl = r.queryBuilder(sl, filter)

	err := sl.ColumnExpr(`SUM("order".total) as total_sales, SUM("order".total_nett) as total_nett, count("oi".id) as total_product`).
		Scan(ctx, &report)
	if err != nil {
		return nil, err
	}

	return &struct {
		TotalSales   float32
		TotalNett    float32
		TotalProduct int
	}{
		TotalSales:   report.TotalSales,
		TotalNett:    report.TotalNett,
		TotalProduct: report.TotalProduct,
	}, nil
}
