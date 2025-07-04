package repository

import (
	"context"

	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/uptrace/bun"
)

type UserStoreAggregateRepository struct {
	db *bun.DB
}

type UserStoreAggregateRepositoryFilter struct {
	UserID  *uint
	StoreID *uint
}

func NewUserStoreAggregateRepository(db *bun.DB) *UserStoreAggregateRepository {
	return &UserStoreAggregateRepository{
		db: db,
	}
}

func (r *UserStoreAggregateRepository) queryBuilder(query *bun.SelectQuery, cond UserStoreAggregateRepositoryFilter) *bun.SelectQuery {
	if cond.UserID != nil {
		query.Where("? = ?", bun.Ident("user_store.user_id"), *cond.UserID)
	}

	if cond.StoreID != nil {
		query.Where("? = ?", bun.Ident("user_store.store_id"), *cond.StoreID)
	}

	return query
}

func (r *UserStoreAggregateRepository) FindOne(ctx context.Context, cond UserStoreAggregateRepositoryFilter) (*models.UserStore, error) {
	query := r.db.NewSelect().Model(&models.UserStore{}).Relation("User").Relation("Store")
	query = r.queryBuilder(query, cond)

	var userStore models.UserStore
	err := query.Scan(ctx, &userStore)
	if err != nil {
		return nil, err
	}

	return &userStore, nil
}
