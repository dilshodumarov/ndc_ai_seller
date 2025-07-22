package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sugurta/internal/pkg/redis"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
	IsAvailable(ctx context.Context) bool
	GetCacheVersion(ctx context.Context, key string) int
	IncrementCacheVersion(ctx context.Context, key string)
}

func NewCache(rdb *redis.RedisDB) *cache {
	return &cache{
		rdb: rdb,
	}
}

type cache struct {
	rdb *redis.RedisDB
}

func (c *cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = c.rdb.Client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		fmt.Println(9999)
		return err
	}
	return nil
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	// tracing
	// ctx, span := otlp_pkg.Start(ctx, "cecheService", "CasheRepoGet")
	// defer span.End()

	data, err := c.rdb.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// span.SetAttributes(
	// 	attribute.Key("data").String(data),
	// )

	return []byte(data), nil
}

func (c *cache) Del(ctx context.Context, key string) error {
	// tracing
	// ctx, span := otlp_pkg.Start(ctx, "cecheService", "CasheRepoDel")
	// defer span.End()

	err := c.rdb.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *cache) IsAvailable(ctx context.Context) bool {
	err := c.rdb.Client.Ping(ctx).Err()
	return err == nil
}

func (c *cache) GetCacheVersion(ctx context.Context, key string) int {
	val, err := c.rdb.Client.Get(ctx, key).Result()
	if err != nil {
		return 0
	}
	v, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return v
}

func (c *cache) IncrementCacheVersion(ctx context.Context, key string) {
	current := c.GetCacheVersion(ctx, key)
	_ = c.rdb.Client.Set(ctx, key, strconv.Itoa(current+1), 0).Err()
}
