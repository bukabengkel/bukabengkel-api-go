package repository

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/uptrace/bun"
)

type StoreRepository struct {
	db *bun.DB
}

type StoreRepositoryFilter struct {
	ID  *uint64
	Key *string
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

