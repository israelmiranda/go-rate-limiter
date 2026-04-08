package ratelimiter

import (
	"context"
	"testing"
	"time"

	"go-rate-limiter/internal/config"
)

type MockStrategy struct {
	data   map[string]int64
	blocks map[string]time.Time
}

func NewMockStrategy() *MockStrategy {
	return &MockStrategy{
		data:   make(map[string]int64),
		blocks: make(map[string]time.Time),
	}
}

func (m *MockStrategy) Increment(ctx context.Context, key string) (int64, error) {
	m.data[key]++
	return m.data[key], nil
}

func (m *MockStrategy) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return time.Second, nil
}

func (m *MockStrategy) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	return nil
}

func (m *MockStrategy) Block(ctx context.Context, key string, duration time.Duration) error {
	m.blocks[key] = time.Now().Add(duration)
	return nil
}

func (m *MockStrategy) IsBlocked(ctx context.Context, key string) (bool, error) {
	if blockTime, exists := m.blocks[key]; exists {
		return time.Now().Before(blockTime), nil
	}
	return false, nil
}

func TestRateLimiter_Allow_IP(t *testing.T) {
	strategy := NewMockStrategy()
	cfg := &config.Config{
		RateLimitIP:   2,
		BlockDuration: time.Minute,
	}
	limiter := NewRateLimiter(strategy, cfg)

	ctx := context.Background()

	// First request should be allowed
	allowed, err := limiter.Allow(ctx, "192.168.1.1", "")
	if err != nil || !allowed {
		t.Errorf("First request should be allowed")
	}

	// Second request should be allowed
	allowed, err = limiter.Allow(ctx, "192.168.1.1", "")
	if err != nil || !allowed {
		t.Errorf("Second request should be allowed")
	}

	// Third request should be blocked
	allowed, err = limiter.Allow(ctx, "192.168.1.1", "")
	if err != nil || allowed {
		t.Errorf("Third request should be blocked")
	}
}

func TestRateLimiter_Allow_Token_Precedence(t *testing.T) {
	strategy := NewMockStrategy()
	cfg := &config.Config{
		RateLimitIP:    1,
		RateLimitToken: 3,
		BlockDuration:  time.Minute,
	}
	limiter := NewRateLimiter(strategy, cfg)

	ctx := context.Background()

	// Token should allow more requests than IP limit
	allowed, err := limiter.Allow(ctx, "192.168.1.1", "token123")
	if err != nil || !allowed {
		t.Errorf("First token request should be allowed")
	}

	allowed, err = limiter.Allow(ctx, "192.168.1.1", "token123")
	if err != nil || !allowed {
		t.Errorf("Second token request should be allowed")
	}

	allowed, err = limiter.Allow(ctx, "192.168.1.1", "token123")
	if err != nil || !allowed {
		t.Errorf("Third token request should be allowed")
	}

	// Fourth should be blocked
	allowed, err = limiter.Allow(ctx, "192.168.1.1", "token123")
	if err != nil || allowed {
		t.Errorf("Fourth token request should be blocked")
	}
}
