package main

import (
	"github.com/gin-gonic/gin"
)

// Handler handles the API request
func Handler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello from Next.go API!",
		"method":  c.Request.Method,
		"path":    c.Request.URL.Path,
	})
}
