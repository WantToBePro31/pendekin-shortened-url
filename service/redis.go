package service

import (
	"context"
	"encoding/json"
	"pendekin/model"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisService interface {
	Get(key string) ([]model.Url, error)
	Store(key string, urls []model.Url) error
	Clear(key string) error
}

type redisService struct {
	rdb *redis.Client
}

func InitRedisService(rdb *redis.Client) *redisService {
	return &redisService{rdb}
}

func (rs *redisService) Get(key string) ([]model.Url, error) {
	cache, err := rs.rdb.Get(context.TODO(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return []model.Url{}, nil
		}
		return []model.Url{}, err
	}

	var urls []model.Url
	if err := json.Unmarshal([]byte(cache), &urls); err != nil {
		return []model.Url{}, err
	}

	return urls, nil
}

func (rs *redisService) Store(key string, urls []model.Url) error {
	urljson, err := json.Marshal(urls)
	if err != nil {
		return err
	}

	err = rs.rdb.Set(context.TODO(), key, urljson, 10*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rs *redisService) Clear(key string) error {
	if err := rs.rdb.Del(context.TODO(), key).Err(); err != nil {
		return err
	}

	return nil
}
