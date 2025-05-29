package repository

import (
	"context"
	"time"

	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/peang/bukabengkel-api-go/src/services/cache_services"
	"github.com/peang/bukabengkel-api-go/src/utils"
	"github.com/uptrace/bun"
)

type OrderDistributorRepository struct {
	db           *bun.DB
	cacheService cache_services.CacheService
}

type OrderDistributorRepositoryFilter struct {
	Key       *string
	StoreID   *uint
	StartDate *time.Time
	EndDate   *time.Time
}

func NewOrderDistributorRepository(db *bun.DB, cacheService cache_services.CacheService) *OrderDistributorRepository {
	return &OrderDistributorRepository{db: db, cacheService: cacheService}
}

func (r *OrderDistributorRepository) queryBuilder(query *bun.SelectQuery, cond OrderDistributorRepositoryFilter) *bun.SelectQuery {
	if cond.Key != nil {
		query.Where("? = ?", bun.Ident("order_distributor.key"), *cond.Key)
	}

	if cond.StoreID != nil {
		query.Where("? = ?", bun.Ident("order_distributor.customer_id"), *cond.StoreID)
	}

	if cond.StartDate != nil && cond.EndDate != nil {
		query.Where("? BETWEEN ? AND ?", bun.Ident("order_distributor.created_at"), cond.StartDate, cond.EndDate)
	}

	return query
}

func (r *OrderDistributorRepository) List(ctx context.Context, page int, perPage int, sort string, cond OrderDistributorRepositoryFilter) ([]models.OrderDistributor, bool, error) {
	sorts := utils.GenerateSort(sort)
	offset, limit := utils.GenerateOffsetLimitV2(page, perPage)

	var hasNext bool = false
	var orderDistributors []models.OrderDistributor

	query := r.db.NewSelect().Column(
		"id",
		"key",
		"distributor_name",
		"total",
		"status",
		"expired_at",
		"paid_at",
		"created_at",
		"updated_at",
	).Model(&orderDistributors)

	query = r.queryBuilder(query, cond)

	err := query.
		Limit(limit).
		Offset(offset).
		OrderExpr(sorts).
		Scan(ctx)
	if err != nil {
		return nil, false, err
	}

	if len(orderDistributors) > perPage {
		orderDistributors = orderDistributors[:len(orderDistributors)-1]
		hasNext = true
	}

	return orderDistributors, hasNext, nil
}

func (r *OrderDistributorRepository) FindOne(ctx context.Context, cond OrderDistributorRepositoryFilter) (*models.OrderDistributor, error) {
	query := r.db.NewSelect().Model(&models.OrderDistributor{})
	query = r.queryBuilder(query, cond)

	var orderDistributor models.OrderDistributor
	err := query.Scan(ctx, &orderDistributor)
	if err != nil {
		return nil, err
	}

	return &orderDistributor, nil
}

func (r *OrderDistributorRepository) CreateOrderDistributor(ctx context.Context, orderDistributor *models.OrderDistributor) error {
	_, err := r.db.NewInsert().Model(orderDistributor).Exec(ctx)
	return err
}

func (r *OrderDistributorRepository) UpdateOrderDistributor(ctx context.Context, orderDistributor *models.OrderDistributor) error {
	_, err := r.db.NewUpdate().Model(orderDistributor).WherePK().Exec(ctx)
	return err
}
