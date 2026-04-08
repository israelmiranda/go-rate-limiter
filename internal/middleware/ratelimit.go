package middleware

import (
	"net/http"

	"go-rate-limiter/internal/ratelimiter"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(limiter *ratelimiter.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		token := c.GetHeader("API_KEY")

		allowed, err := limiter.Allow(c.Request.Context(), ip, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				"you have reached the maximum number of requests or actions allowed within a certain time frame",
			)
			return
		}

		c.Next()
	}
}
