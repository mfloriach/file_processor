package utils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

var c *cache

type cache struct {
	rdb *redis.Client
}

func GetRedisClient() *cache {
	if c != nil {
		return c
	}

	env := GetEnv()

	return &cache{
		rdb: redis.NewClient(&redis.Options{
			Addr:     env.RedisAddr,
			Password: env.RedisPassword,
			DB:       0,
		}),
	}
}

func (c cache) Set(ctx context.Context, key string, value interface{}, expiration ...int64) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := c.rdb.Set(ctx, key, b, time.Minute); err.Err() != nil {
		return err.Err()
	}

	return nil
}
