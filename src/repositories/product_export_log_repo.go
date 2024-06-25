package repository

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/peang/bukabengkel-api-go/src/utils"
	"github.com/uptrace/bun"
)

type ProductExportLogRepository struct {
	db *bun.DB
}

type ProductExportLogRepositoryFilter struct {
	Name    *string
	StoreID *uint
	UserID  *uint
}

func NewProductExportLogRepository(db *bun.DB) *ProductExportLogRepository {
	return &ProductExportLogRepository{
		db: db,
	}
}

func (r *ProductExportLogRepository) queryBuilder(query *bun.SelectQuery, cond ProductExportLogRepositoryFilter) *bun.SelectQuery {
	if cond.StoreID != nil {
		query.Where("? = ?", bun.Ident("product.store_id"), cond.StoreID)
	}

	if cond.UserID != nil {
		query.Where("? > ?", bun.Ident("product.user_id"), cond.UserID)
	}

	return query
}

func (r *ProductExportLogRepository) List(ctx context.Context, page int, perPage int, sort string, filter ProductExportLogRepositoryFilter) (*[]models.ProductExportLog, int, error) {
	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimit(page, perPage)

	var productsExportLog []models.ProductExportLog
	sl := r.db.NewSelect().Model(&productsExportLog)
	sl = r.queryBuilder(sl, filter)

	count, err := sl.
		Relation("Store").
		Relation("User").
		Limit(limit).Offset(offset).OrderExpr(sorts).ScanAndCount(context.TODO())
	if err != nil {
		return nil, 0, err
	}

	if len(productsExportLog) == 0 {
		return &[]models.ProductExportLog{}, count, nil
	}

	var entityProductExportLogs []models.ProductExportLog
	entityProductExportLogs = append(entityProductExportLogs, productsExportLog...)

	return &entityProductExportLogs, count, nil
}
