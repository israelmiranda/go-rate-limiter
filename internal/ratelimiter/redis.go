package ratelimiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStrategy struct {
	client *redis.Client
}

func NewRedisStrategy(addr, password string, db int) *RedisStrategy {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStrategy{client: rdb}
}

func (r *RedisStrategy) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *RedisStrategy) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return ttl, nil
}

func (r *RedisStrategy) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}

func (r *RedisStrategy) Block(ctx context.Context, key string, duration time.Duration) error {
	blockKey := "block:" + key
	return r.client.Set(ctx, blockKey, "1", duration).Err()
}

func (r *RedisStrategy) IsBlocked(ctx context.Context, key string) (bool, error) {
	blockKey := "block:" + key
	exists, err := r.client.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
