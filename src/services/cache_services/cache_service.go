package cache_services

import (
	"context"
	"errors"
	"time"

	"github.com/peang/bukabengkel-api-go/src/config"
)

type CacheService interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

const (
	CACHE_SERVICE_REDIS = "redis"
)

func NewCacheService(config *config.Config) (CacheService, error) {
	switch config.Cache.CacheServiceName {
	case CACHE_SERVICE_REDIS:
		return newRedisCache(config), nil
	default:
		return nil, errors.New("cache service not found")
	}
}
