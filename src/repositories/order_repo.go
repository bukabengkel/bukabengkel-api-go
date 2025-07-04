package repository

import (
	"context"
	"time"

	"github.com/peang/bukabengkel-api-go/src/services/cache_services"
	"github.com/peang/bukabengkel-api-go/src/utils"
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

type productOrderResult struct {
	ProductKey      string
	ProductName     string
	ProductCategory string
	ProductUnit     string
	QtySales        int
	QtyStock        float64
	OrderDate       string
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
		query.Where("? = ?", bun.Ident("order.store_id"), *cond.StoreID)
	}

	if cond.StartDate != nil && cond.EndDate != nil {
		query.Where("? BETWEEN ? AND ?", bun.Ident("order.order_date"), cond.StartDate, cond.EndDate)
	}

	return query
}

func (r *OrderRepository) OrderSalesReport(ctx context.Context, filter OrderRepositoryFilter) (*salesOrderResult, error) {
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

func (r *OrderRepository) ProductSalesReport(ctx context.Context, page int, perPage int, filter OrderRepositoryFilter) (*[]productOrderResult, *int, error) {
	offset, limit := utils.GenerateOffsetLimit(page, perPage)

	var results []productOrderResult

	sl := r.db.NewSelect().Table(`order`).
		Join(`INNER JOIN "order_item" ON order_item.order_id = "order".id`).
		Join(`INNER JOIN "product" ON order_item.product_key_id::uuid = product.key`).
		Join(`INNER JOIN "product_category" ON product.category_id = product_category.id`)

	sl = r.queryBuilder(sl, filter)

	count, err := sl.ColumnExpr(`product.key as product_key, product.name as product_name, product_category.name as product_category, product.unit as product_unit, SUM(order_item.qty) as qty_sales, product.stock as qty_stock, MAX("order".order_date) AS order_date`).
		GroupExpr(`product.key, product.name, product_category.name, product.unit, product.stock`).
		OrderExpr(`order_date DESC`).
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx, &results)

	if err != nil {
		return nil, nil, err
	}

	return &results, &count, nil
}

func (r *OrderRepository) Count(ctx context.Context, filter OrderRepositoryFilter) (int, error) {
	sl := r.db.NewSelect().Table("order")
	sl = r.queryBuilder(sl, filter)

	count, err := sl.Count(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// func GenerateHashKey(filter OrderRepositoryFilter) string {
// 	id := fmt.Sprint("report_sales_", filter.StoreID, "_", filter.StartDate, "_", filter.EndDate)

// 	hash := md5.New()
// 	hash.Write([]byte(id))
// 	hashedBytes := hash.Sum(nil)
// 	hashedString := hex.EncodeToString(hashedBytes)

// 	return hashedString
// }
