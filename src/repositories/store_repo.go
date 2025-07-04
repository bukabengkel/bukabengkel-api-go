package repository

import (
	"context"
	"time"

	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/uptrace/bun"
)

type StoreRepository struct {
	db *bun.DB
}

type StoreRepositoryFilter struct {
	ID  *uint64
	Key *string
	StartDate *time.Time
	EndDate *time.Time
}

func NewStoreRepository(db *bun.DB) *StoreRepository {
	return &StoreRepository{
		db: db,
	}
}

func (r *StoreRepository) queryBuilder(query *bun.SelectQuery, cond StoreRepositoryFilter) *bun.SelectQuery {
	if cond.ID != nil {
		query.Where("? = ?", bun.Ident("store.id"), *cond.ID)
	}

	if cond.Key != nil {
		query.Where("? = ?", bun.Ident("store.key"), *cond.Key)
	}

	if cond.StartDate != nil && cond.EndDate != nil {
		query.Where("? BETWEEN ? AND ?", bun.Ident("store.created_at"), cond.StartDate, cond.EndDate)
	}

	return query
}

func (r *StoreRepository) FindOne(ctx context.Context, cond StoreRepositoryFilter) (*models.Store, error) {
	query := r.db.NewSelect().Model(&models.Store{})
	query = r.queryBuilder(query, cond)

	var store models.Store
	err := query.Scan(ctx, &store)
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (r *StoreRepository) Count(ctx context.Context, filter StoreRepositoryFilter) (int, error) {
	sl := r.db.NewSelect().Table("store")
	sl = r.queryBuilder(sl, filter)

	count, err := sl.Count(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}