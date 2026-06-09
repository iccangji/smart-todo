package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string, dest any) bool
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
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

func (r *RedisCache) Get(ctx context.Context, key string, dest any) bool {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return false
	}

	return true
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	r.client.Set(ctx, key, data, ttl)
	return nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) {
	r.client.Del(ctx, key)
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
