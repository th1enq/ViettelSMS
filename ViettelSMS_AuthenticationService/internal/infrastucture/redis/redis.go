package rdb

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/th1enq/ViettelSMS_AuthenticationService/internal/config"
)

type redisDB struct {
	client *redis.Client
}

func NewRedisDB(cfg *config.Config) (CacheEngine, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	fmt.Printf("Redis Config: %+v\n", cfg.Redis)

	ctx := context.Background()
	_, err := r.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return &redisDB{client: r}, nil
}

func (r *redisDB) GetCache() *redis.Client {
	return r.client
}
