package repository

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/peang/bukabengkel-api-go/src/services/cache_services"
	"github.com/uptrace/bun"
)

type OrderRepository struct {
	db           *bun.DB
	cacheService cache_services.CacheService
}

type OrderRepositoryFilter struct {
	StoreID   *uint
	StartDate *time.Time
	EndDate   *time.Time
}

type salesOrderResult struct {
	TotalSales   float32
	TotalNett    float32
	TotalProduct int
}

func NewOrderRepository(
	db *bun.DB,
	cacheService cache_services.CacheService,
) *OrderRepository {
	return &OrderRepository{
		db:           db,
		cacheService: cacheService,
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

func (r *OrderRepository) CountReport(ctx context.Context, filter OrderRepositoryFilter) (*salesOrderResult, error) {
	cacheKey := generateHashKey(filter)
	fmt.Println(cacheKey)
	cache, err := r.cacheService.Get(ctx, cacheKey)
	if err == nil {
		return nil, err
	}
	if cache != nil {
		return cache.(*salesOrderResult), nil
	}

	var report struct {
		TotalSales   float32
		TotalNett    float32
		TotalProduct int
	}
	sl := r.db.NewSelect().Table("order").Join(`LEFT JOIN order_item as oi ON "oi".id = "order".id`)
	sl = r.queryBuilder(sl, filter)

	err = sl.ColumnExpr(`SUM("order".total) as total_sales, SUM("order".total_nett) as total_nett, count("oi".id) as total_product`).
		Scan(ctx, &report)
	if err != nil {
		return nil, err
	}

	result := salesOrderResult{
		TotalSales:   report.TotalSales,
		TotalNett:    report.TotalNett,
		TotalProduct: report.TotalProduct,
	}

	// err = r.cacheService.Set(ctx, cacheKey, result)
	// if err != nil {
	// 	return nil, err
	// }

	return &result, nil
}

func generateHashKey(filter OrderRepositoryFilter) string {
	id := fmt.Sprint("report_sales_", filter.StoreID, "_", filter.StartDate, "_", filter.EndDate)

	hash := md5.New()
	hash.Write([]byte(id))
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}
