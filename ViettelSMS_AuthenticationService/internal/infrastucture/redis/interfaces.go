package rdb

import "github.com/redis/go-redis/v9"

type CacheEngine interface {
	GetCache() *redis.Client
}
