package redisfacade

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisFacade struct {
	RedisClient           *redis.Client
	CacheInvalidationTime time.Duration
	CacheSearchRadius     float64
}

func NewRedisFacade(cli *redis.Client, cit time.Duration, csr float64) *RedisFacade {
	return &RedisFacade{
		RedisClient:           cli,
		CacheInvalidationTime: cit,
		CacheSearchRadius:     csr,
	}
}
