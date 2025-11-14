package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	lastSeen time.Time
	count    int
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(requestsPerWindow int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    requestsPerWindow,
		window:   window,
	}
	
	// Cleanup old visitors periodically
	go rl.cleanupVisitors()
	
	return rl
}

func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		rl.mu.Lock()
		defer rl.mu.Unlock()
		
		v, exists := rl.visitors[ip]
		now := time.Now()
		
		if !exists {
			rl.visitors[ip] = &visitor{lastSeen: now, count: 1}
			c.Next()
			return
		}
		
		// Reset counter if window has passed
		if now.Sub(v.lastSeen) > rl.window {
			v.count = 1
			v.lastSeen = now
			c.Next()
			return
		}
		
		// Check if limit exceeded
		if v.count >= rl.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}
		
		v.count++
		v.lastSeen = now
		c.Next()
	}
}
