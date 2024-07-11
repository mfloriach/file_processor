package utils

import (
	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func GetRedisClient() *redis.Client {
	if rdb != nil {
		return rdb
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81", // no password set
		DB:       0,                                  // use default DB
	})

	return rdb
}
