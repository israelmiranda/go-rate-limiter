package ratelimiter

import (
	"context"
	"time"
)

type PersistenceStrategy interface {
	Increment(ctx context.Context, key string) (int64, error)
	GetTTL(ctx context.Context, key string) (time.Duration, error)
	SetTTL(ctx context.Context, key string, ttl time.Duration) error
	Block(ctx context.Context, key string, duration time.Duration) error
	IsBlocked(ctx context.Context, key string) (bool, error)
}
