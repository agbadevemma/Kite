package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)



type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)


func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	c, exists := clients[ip]

	if !exists {
		limiter := rate.NewLimiter(2, 5) 


		clients[ip] = &client{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	c.lastSeen = time.Now()
	return c.limiter
}


func RateLimitMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		limiter := getLimiter(ip)

		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		ctx.Next()
	}
}