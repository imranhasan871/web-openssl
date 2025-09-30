package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Clean up old entries every minute
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *RateLimiter) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	for key, timestamps := range rl.requests {
		// Remove timestamps older than the window
		filtered := make([]time.Time, 0)
		for _, ts := range timestamps {
			if now.Sub(ts) <= rl.window {
				filtered = append(filtered, ts)
			}
		}

		if len(filtered) == 0 {
			delete(rl.requests, key)
		} else {
			rl.requests[key] = filtered
		}
	}
}

func (rl *RateLimiter) isAllowed(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	timestamps := rl.requests[key]

	// Remove timestamps older than the window
	filtered := make([]time.Time, 0)
	for _, ts := range timestamps {
		if now.Sub(ts) <= rl.window {
			filtered = append(filtered, ts)
		}
	}

	// Check if we're within the limit
	if len(filtered) >= rl.limit {
		return false
	}

	// Add current request
	filtered = append(filtered, now)
	rl.requests[key] = filtered

	return true
}

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		// Use IP address as the key for rate limiting
		key := c.ClientIP()

		// For authenticated users, use user ID instead
		if userID, exists := c.Get("user_id"); exists {
			key = fmt.Sprintf("user_%v", userID)
		}

		if !limiter.isAllowed(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": int(window.Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}