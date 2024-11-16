package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(config *RateLimiterConfig) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimiter(config))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return router
}

func TestRateLimiter(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		requests   int
		capacity   float64
		refillRate float64
		sleep      time.Duration
		wantCodes  []int
	}{
		{
			name:       "Basic rate limiting",
			requests:   5,
			capacity:   3,
			refillRate: 1,
			wantCodes:  []int{200, 200, 200, 429, 429},
		},
		{
			name:       "With refill",
			requests:   4,
			capacity:   1,
			refillRate: 1,
			sleep:      time.Second,
			wantCodes:  []int{200, 429, 200, 429}, // Use initial token, fail, refill and use token, fail
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			router := gin.New()

			config := &RateLimiterConfig{
				Capacity:   tt.capacity,
				RefillRate: tt.refillRate,
			}

			router.Use(RateLimiter(config))

			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			for i := 0; i < tt.requests; i++ {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/test", nil)
				router.ServeHTTP(w, req)

				assert.Equal(t, tt.wantCodes[i], w.Code, "Request %d", i)

				if tt.sleep > 0 && i%2 == 1 { // Sleep after every second request
					time.Sleep(tt.sleep)
				}
			}
		})
	}
}

func TestRateLimiterCustomClient(t *testing.T) {
	t.Parallel()

	config := &RateLimiterConfig{
		Capacity:   2,
		RefillRate: 1,
		Client: func(c *gin.Context) string {
			return c.GetHeader("X-API-Key")
		},
	}

	router := setupRouter(config)
	clients := []string{"client1", "client2"}

	for _, client := range clients {
		client := client
		t.Run("Client "+client, func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 3; i++ {
				w := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("X-API-Key", client)
				router.ServeHTTP(w, req)

				expectedCode := http.StatusOK
				if i >= 2 {
					expectedCode = http.StatusTooManyRequests
				}
				assert.Equal(t, expectedCode, w.Code,
					"Request %d for client %s failed: expected status %d, got %d",
					i+1, client, expectedCode, w.Code)
			}
		})
	}
}

func TestRateLimiterConcurrent(t *testing.T) {
	t.Parallel()

	config := &RateLimiterConfig{
		Capacity:   5,
		RefillRate: 1,
	}

	router := setupRouter(config)
	const concurrentRequests = 10
	results := make(chan int, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)
			results <- w.Code
		}()
	}

	successCount := 0
	rateLimitCount := 0

	for i := 0; i < concurrentRequests; i++ {
		code := <-results
		switch code {
		case http.StatusOK:
			successCount++
		case http.StatusTooManyRequests:
			rateLimitCount++
		}
	}

	assert.Equal(t, concurrentRequests, successCount+rateLimitCount,
		"Total responses should equal total requests")
	assert.Equal(t, 5, successCount,
		"Should allow exactly capacity (5) successful requests")
	assert.Equal(t, 5, rateLimitCount,
		"Should rate limit remaining requests")
}
