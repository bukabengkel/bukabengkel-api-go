package repository

import (
	"context"
	"database/sql"

	"github.com/peang/bukabengkel-api-go/src/models"
	"github.com/uptrace/bun"
)

type LocationRepository struct {
	db *bun.DB
}

type LocationRepositoryFilter struct {
	EntityID   *uint64
	EntityType *string
}

func NewLocationRepository(db *bun.DB) *LocationRepository {
	return &LocationRepository{
		db: db,
	}
}

func (r *LocationRepository) queryBuilder(query *bun.SelectQuery, cond LocationRepositoryFilter) *bun.SelectQuery {
	if cond.EntityID != nil {
		query.Where("? = ?", bun.Ident("entity_id"), cond.EntityID)
	}

	if cond.EntityType != nil {
		query.Where("? = ?", bun.Ident("entity_type"), cond.EntityType)
	}

	return query
}

func (r *LocationRepository) FindOne(ctx context.Context, cond LocationRepositoryFilter) (*models.RajaOngkirLocation, error) {
	query := r.db.NewSelect().Model(&models.RajaOngkirLocation{})
	query = r.queryBuilder(query, cond)

	var location models.RajaOngkirLocation
	err := query.Scan(ctx, &location)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &location, nil
}
