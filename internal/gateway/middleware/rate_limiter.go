package middleware

import (
	"net/http"
	"sync"
	"time"
	
	"github.com/AhmedHossam777/task-orchestrator/pkg/response"
	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64
	lastRefill time.Time
}

func (rl *rateLimiter) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.tokens += elapsed * rl.refillRate
	if rl.tokens > rl.maxTokens {
		rl.tokens = rl.maxTokens
	}
	rl.lastRefill = now
	
	// try to consume
	if rl.tokens >= 1.0 {
		rl.tokens--
		return true
	}
	
	return false
}

func RateLimiter(maxRequest float64, perSecond float64) gin.HandlerFunc {
	clients := &sync.Map{}
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		value, _ := clients.LoadOrStore(
			ip, &rateLimiter{
				tokens:     maxRequest,
				maxTokens:  maxRequest,
				refillRate: perSecond,
				lastRefill: time.Now(),
			},
		)
		
		limiter := value.(*rateLimiter)
		if !limiter.allow() {
			response.Fail(
				c, http.StatusTooManyRequests,
				"RATE_LIMIT_EXCEEDED",
				"too many requests, please try again later",
			)
			c.Abort()
			return
		}
		
		c.Next()
	}
}
