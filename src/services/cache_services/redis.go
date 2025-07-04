package cache_services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/peang/bukabengkel-api-go/src/config"
	"github.com/peang/bukabengkel-api-go/src/models"
)

type redisCache struct {
	client *redis.Client
}

func newRedisCache(config *config.Config) CacheService {
	return &redisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.Cache.CacheHost, config.Cache.CachePort),
			Username: config.Cache.CacheUsername,
			Password: config.Cache.CachePassword,
			DB:       0,
		}),
	}
}

func (r *redisCache) Get(ctx context.Context, key string) (value any, err error) {
	result, err := r.client.Get(ctx, key).Result()
	fmt.Println("result", result)
	if result == "" {
		return nil, err
	}

	return result, err
}

func (r *redisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	// hashValue, err := HashValue(value)
	// if err != nil {
	// 	return err
	// }

	if ttl == 0 {
		ttl = time.Hour * 1
	}

	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func HashValue(value any) (string, error) {
	by, err := convertToBytes(value)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	hash.Write(by)
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString, nil
}

// convertToBytes converts the input value to a byte slice
func convertToBytes(input any) ([]byte, error) {
	switch v := input.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	case *[]models.Product:
		return json.Marshal(v)
	default:
		// For other types, use reflection to convert to bytes
		valueType := reflect.TypeOf(input)
		value := reflect.ValueOf(input)

		if valueType.Kind() == reflect.Ptr {
			valueType = valueType.Elem()
			value = value.Elem()
		}

		if valueType.Kind() == reflect.Slice {
			return value.Bytes(), nil
		}

		// For other types, attempt to convert to string representation
		str := fmt.Sprintf("%v", input)
		return []byte(str), nil
	}
}
