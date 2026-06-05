package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]string, bool)
	Set(ctx context.Context, key string, value []string, ttl time.Duration)
	Delete(ctx context.Context, key string)
}

type RedisCache struct {
	client *redis.Client
}

func ConnectRedis(ctx context.Context, addr string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	return &RedisCache{
		client: rdb,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]string, bool) {

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, false
	}

	var data []string
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, false
	}

	return data, true
}

func (r *RedisCache) Set(ctx context.Context, key string, value []string, ttl time.Duration) {

	bytes, _ := json.Marshal(value)

	r.client.Set(ctx, key, bytes, ttl)
}

func (r *RedisCache) Delete(ctx context.Context, key string) {
	r.client.Del(ctx, key)
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
