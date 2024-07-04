package redis

import (
	"context"
	"github.com/online-store/pkg/cache"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

type redisRepository struct {
	redisClient *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) cache.RedisRepository {
	return &redisRepository{redisClient: redisClient}
}

func (r redisRepository) Fetch(ctx context.Context, key string) (*string, error) {
	result, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r redisRepository) Save(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	jsonString, err := jsoniter.MarshalToString(data)
	if err != nil {
		return err
	}

	if err := r.redisClient.Set(ctx, key, jsonString, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func (r redisRepository) Delete(ctx context.Context, key string) error {
	var err error
	var iterator = r.redisClient.Scan(ctx, 0, key+"*", 0).Iterator()
	for iterator.Next(ctx) {
		err = r.redisClient.Del(ctx, iterator.Val()).Err()
	}
	if err := iterator.Err(); err != nil {
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

func (r redisRepository) Deletes(ctx context.Context, key []string) error {
	var err error
	for _, v := range key {
		var iterator = r.redisClient.Scan(ctx, 0, v+"*", 0).Iterator()
		if err := iterator.Err(); err != nil {
			return err
		}
		for iterator.Next(ctx) {
			err = r.redisClient.Del(ctx, iterator.Val()).Err()
		}

		if !iterator.Next(ctx) {
			err = r.redisClient.Del(ctx, v).Err()
		}
		if err != nil {
			return err
		}
	}
	return nil
}
