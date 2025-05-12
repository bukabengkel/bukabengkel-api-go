package config

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func LoadRedis(c *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Cache.CacheHost, c.Cache.CachePort),
		Username: c.Cache.CacheUsername,
		Password: c.Cache.CachePassword,
		DB:       1,
	})
}