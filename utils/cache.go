package utils

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

func InitRedis() (*redis.Client, error) {
	rdc := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PWD"),
		DB:       0,
	})

	_, err := rdc.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	return rdc, nil
}