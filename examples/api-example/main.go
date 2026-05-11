package main

import (
	"fmt"
	
	"github.com/gin-gonic/gin"
)

// This is an example of an API handler in Next.go
// File: app/api/users/handler.go

func GetUsers(c *gin.Context) {
	users := []map[string]interface{}{
		{"id": 1, "name": "John Doe"},
		{"id": 2, "name": "Jane Smith"},
	}
	c.JSON(200, gin.H{
		"users": users,
	})
}

func CreateUser(c *gin.Context) {
	var user map[string]interface{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(201, gin.H{
		"message": "User created",
		"user":    user,
	})
}

func Handler(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		GetUsers(c)
	case "POST":
		CreateUser(c)
	default:
		c.JSON(405, gin.H{"error": "Method not allowed"})
	}
}

func main() {
	// This is just an example - in Next.go, handlers are auto-discovered
	r := gin.Default()
	r.GET("/api/users", GetUsers)
	r.POST("/api/users", CreateUser)
	
	fmt.Println("Example API server running on :8080")
	r.Run(":8080")
}
