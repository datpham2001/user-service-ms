package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/datpham/user-service-ms/config"
	"github.com/datpham/user-service-ms/pkg/logger"
	"github.com/redis/go-redis/v9"
)

const (
	ServiceCachePrefix = "user-service-ms"
)

type Cache struct {
	logger      *logger.Logger
	cacheClient *redis.Client
}

func NewCacheClient(logger *logger.Logger, cfg *config.Config) *Cache {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Cache.Host, cfg.Cache.Port),
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DB,
	})

	return &Cache{
		logger:      logger,
		cacheClient: redisClient,
	}
}

func (c *Cache) Ping(ctx context.Context) error {
	if err := c.cacheClient.Ping(ctx).Err(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) Get(ctx context.Context, key string, obj any) error {
	key = fmt.Sprintf("%s:%s", ServiceCachePrefix, key)
	result, err := c.cacheClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(result), &obj)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	key = fmt.Sprintf("%s:%s", ServiceCachePrefix, key)
	if err := c.cacheClient.Set(ctx, key, value, expiration).Err(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	key = fmt.Sprintf("%s:%s", ServiceCachePrefix, key)
	if err := c.cacheClient.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	key = fmt.Sprintf("%s:%s", ServiceCachePrefix, key)
	return c.cacheClient.TTL(ctx, key).Result()
}

func (c *Cache) Incr(ctx context.Context, key string) (int64, error) {
	key = fmt.Sprintf("%s:%s", ServiceCachePrefix, key)
	return c.cacheClient.Incr(ctx, key).Result()
}

func (c *Cache) Close() error {
	return c.cacheClient.Close()
}
