package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-rate-limiter/internal/config"
	"go-rate-limiter/internal/middleware"
	"go-rate-limiter/internal/ratelimiter"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	strategy := ratelimiter.NewRedisStrategy(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)

	limiter := ratelimiter.NewRateLimiter(strategy, cfg)

	router := setupRouter(limiter)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("\nShutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited properly.")
}

func setupRouter(limiter *ratelimiter.RateLimiter) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(middleware.RateLimitMiddleware(limiter))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Request allowed"})
	})

	return router
}
