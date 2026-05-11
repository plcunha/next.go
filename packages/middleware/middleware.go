package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Middleware represents a Next.go middleware
type Middleware interface {
	Handle(c *gin.Context)
}

// MiddlewareFunc is a function type that implements Middleware
type MiddlewareFunc func(c *gin.Context)

func (f MiddlewareFunc) Handle(c *gin.Context) {
	f(c)
}

// Chain represents a chain of middlewares
type Chain struct {
	middlewares []Middleware
}

// NewChain creates a new middleware chain
func NewChain() *Chain {
	return &Chain{
		middlewares: make([]Middleware, 0),
	}
}

// Use adds a middleware to the chain
func (c *Chain) Use(m Middleware) *Chain {
	c.middlewares = append(c.middlewares, m)
	return c
}

// UseFunc adds a middleware function to the chain
func (c *Chain) UseFunc(fn func(*gin.Context)) *Chain {
	c.middlewares = append(c.middlewares, MiddlewareFunc(fn))
	return c
}

// Handler returns a gin middleware handler
func (c *Chain) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, m := range c.middlewares {
			m.Handle(ctx)
			if ctx.IsAborted() {
				return
			}
		}
	}
}

// Common middlewares

// Logger returns a logger middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fmt.Printf("[NEXT-GO] %s %s %d %v\n", 
			c.Request.Method, path, status, latency)
	}
}

// CORS returns a CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", 
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", 
			"POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Security returns a security headers middleware
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		c.Next()
	}
}

// Compression returns a compression middleware
func Compression() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple compression - in production use proper gzip
		c.Next()
	}
}

// RateLimit returns a rate limiting middleware
func RateLimit(maxRequests int, duration time.Duration) gin.HandlerFunc {
	requests := make(map[string][]time.Time)
	
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		
		// Clean old requests
		if reqs, ok := requests[ip]; ok {
			var valid []time.Time
			for _, t := range reqs {
				if now.Sub(t) < duration {
					valid = append(valid, t)
				}
			}
			requests[ip] = valid
		}
		
		// Check limit
		if len(requests[ip]) >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		
		requests[ip] = append(requests[ip], now)
		c.Next()
	}
}

// Auth middleware for authentication
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for auth token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
		
		// Validate token (simplified)
		if len(token) < 10 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
