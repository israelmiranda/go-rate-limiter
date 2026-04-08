package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"go-rate-limiter/internal/config"
)

type RateLimiter struct {
	strategy      PersistenceStrategy
	ipLimit       int
	tokenLimit    int
	blockDuration time.Duration
}

func NewRateLimiter(strategy PersistenceStrategy, cfg *config.Config) *RateLimiter {
	return &RateLimiter{
		strategy:      strategy,
		ipLimit:       cfg.RateLimitIP,
		tokenLimit:    cfg.RateLimitToken,
		blockDuration: cfg.BlockDuration,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, ip, token string) (bool, error) {
	var key string
	var limit int

	if token != "" {
		key = fmt.Sprintf("token:%s", token)
		limit = rl.tokenLimit
	} else {
		key = fmt.Sprintf("ip:%s", ip)
		limit = rl.ipLimit
	}

	blocked, err := rl.strategy.IsBlocked(ctx, key)
	if err != nil {
		return false, err
	}
	if blocked {
		return false, nil
	}

	count, err := rl.strategy.Increment(ctx, key)
	if err != nil {
		return false, err
	}

	if count == 1 {
		err = rl.strategy.SetTTL(ctx, key, time.Second)
		if err != nil {
			return false, err
		}
	}

	if count > int64(limit) {
		err = rl.strategy.Block(ctx, key, rl.blockDuration)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
