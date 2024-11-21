package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Limiter represents a token bucket rate limiter
type Limiter struct {
	tokens        float64
	capacity      float64
	refillRate    float64
	lastTimestamp time.Time
	mu            sync.Mutex
}

type RateLimiterConfig struct {
	Capacity   float64                   // Maximum tokens
	RefillRate float64                   // Tokens per second
	Client     func(*gin.Context) string // Function to identify clients
}

func NewRateLimiter(cfg *RateLimiterConfig) *Limiter {
	return &Limiter{
		tokens:        cfg.Capacity,
		capacity:      cfg.Capacity,
		refillRate:    cfg.RefillRate,
		lastTimestamp: time.Now(),
	}
}

func (rl *Limiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastTimestamp)
	tokensToAdd := elapsed.Seconds() * rl.refillRate

	if tokensToAdd > 0 {
		rl.tokens = min(rl.capacity, rl.tokens+tokensToAdd)
		rl.lastTimestamp = now
	}
}

func (rl *Limiter) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()

	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}
	return false
}

func RateLimiter(cfg *RateLimiterConfig) gin.HandlerFunc {
	if cfg == nil {
		log.Fatal("RateLimiter Config is nil!")
	}

	limiters := &sync.Map{}

	if cfg.Client == nil {
		cfg.Client = func(c *gin.Context) string {
			return c.ClientIP()
		}
	}

	return func(c *gin.Context) {
		clientID := cfg.Client(c)

		limiterI, _ := limiters.LoadOrStore(clientID, NewRateLimiter(cfg))
		limiter := limiterI.(*Limiter)

		if !limiter.allow() {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%.0f", cfg.Capacity))

			retryAfter := time.Duration(1000/cfg.RefillRate) * time.Millisecond
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%.0f", retryAfter.Seconds()))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": retryAfter.String(),
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
