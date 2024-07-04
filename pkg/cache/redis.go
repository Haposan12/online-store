package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v8"
	"strconv"
	"time"
)

type RedisConfig struct {
	Key           string `json:"key"`
	RedisURL      string `json:"conn"`
	RedisPassword string `json:"password"`
	RedisDatabase string `json:"dbNum"`
}

func NewRedisClient(cfg RedisConfig) *redis.Client {
	redisDb, err := strconv.Atoi(cfg.RedisDatabase)
	if err != nil {
		redisDb = 1 // default db
	}
	var options = &redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       redisDb,
	}
	redisClient := redis.NewClient(options)
	redisClient.AddHook(nrredis.NewHook(options))
	return redisClient
}

func SetConf(conf string) RedisConfig {
	var model RedisConfig
	json.Unmarshal([]byte(conf), &model)
	return model
}

type RedisRepository interface {
	Fetch(ctx context.Context, key string) (*string, error)
	Save(ctx context.Context, key string, data interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Deletes(ctx context.Context, key []string) error
}
