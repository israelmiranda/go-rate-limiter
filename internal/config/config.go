package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	RateLimitIP    int
	RateLimitToken int
	BlockDuration  time.Duration

	ServerPort string
}

func Load() *Config {
	cfg := &Config{}

	cfg.RedisAddr = getEnv("REDIS_ADDR", "localhost:6379")
	cfg.RedisPassword = getEnv("REDIS_PASSWORD", "")
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	cfg.RedisDB = redisDB

	rateLimitIP, _ := strconv.Atoi(getEnv("RATE_LIMIT_IP", "10"))
	cfg.RateLimitIP = rateLimitIP

	rateLimitToken, _ := strconv.Atoi(getEnv("RATE_LIMIT_TOKEN", "100"))
	cfg.RateLimitToken = rateLimitToken

	blockDuration, _ := time.ParseDuration(getEnv("BLOCK_DURATION", "5m"))
	cfg.BlockDuration = blockDuration

	cfg.ServerPort = getEnv("SERVER_PORT", "8080")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
