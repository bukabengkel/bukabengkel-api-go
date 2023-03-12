package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/peang/bukabengkel-api-go/src/domain/entity"
	repo "github.com/peang/bukabengkel-api-go/src/domain/repositories"
	"github.com/peang/bukabengkel-api-go/src/utils"
)

type productRepository struct {
	db *pgx.Conn
}

func NewProductRepository(db *pgx.Conn) repo.ProductRepositoryInterface {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) List(ctx context.Context, page int, perPage int, sort string, filter entity.ProductEntityRepositoryFilter) (products []*entity.ProductEntity, count int, err error) {
	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimit(page, perPage)

	query := `
		SELECT 
			p.id, p.key, p.category_id, pc.name, p.name, p.slug, p.description, p.unit, p.price, p.sell_price, p.stock, p.stock_minimum, p.is_stock_unlimited, p.status 
		FROM
			product p
			LEFT JOIN product_category pc ON p.category_id = pc.id`
	conditions := make([]string, 0)

	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ilike '%%%s%%'", filter.Name))
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sorts.Field, strings.ToUpper(sorts.Method))
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	var rows pgx.Rows
	rows, err = r.db.Query(ctx, query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		product := &entity.ProductEntity{}
		err = rows.Scan(
			&product.ID,
			&product.Key,
			&product.Category.ID,
			&product.Category.Name,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Unit,
			&product.Price,
			&product.SellPrice,
			&product.Stock,
			&product.StockMinimum,
			&product.IsStockUnlimited,
			&product.Status,
		)
		if err != nil {
			return
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return
	}

	countQuery := "SELECT COUNT(*) FROM product"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	err = r.db.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		return
	}

	return
}
