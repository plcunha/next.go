package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler represents an API handler
type Handler struct {
	Path       string
	Methods    []string
	Handler    gin.HandlerFunc
}

// RouteGroup represents a group of API routes
type RouteGroup struct {
	Prefix string
	Routes []Handler
}

// RegisterAPIRoutes registers all API routes from the app/api directory
func RegisterAPIRoutes(r *gin.Engine, appDir string) error {
	apiDir := filepath.Join(appDir, "app", "api")

	if _, err := os.Stat(apiDir); os.IsNotExist(err) {
		return nil // No API directory
	}

	// Walk through API directory
	err := filepath.Walk(apiDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Process handler files
		if !isHandlerFile(path) {
			return nil
		}

		return registerAPIFile(r, apiDir, path)
	})

	return err
}

// registerAPIFile registers routes from an API handler file
func registerAPIFile(r *gin.Engine, apiDir, filePath string) error {
	relPath, _ := filepath.Rel(apiDir, filePath)
	urlPath := fileToURL(relPath)

	// Get handler functions
	handlers := getHandlersFromFile(filePath)

	// Register routes
	for method, handler := range handlers {
		switch method {
		case "GET":
			r.GET(urlPath, handler)
		case "POST":
			r.POST(urlPath, handler)
		case "PUT":
			r.PUT(urlPath, handler)
		case "DELETE":
			r.DELETE(urlPath, handler)
		case "PATCH":
			r.PATCH(urlPath, handler)
		}
	}

	fmt.Printf("  ✓ API: %s -> %s\n", urlPath, filePath)
	return nil
}

// getHandlersFromFile parses a handler file and returns HTTP method handlers
func getHandlersFromFile(filePath string) map[string]gin.HandlerFunc {
	handlers := make(map[string]gin.HandlerFunc)

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return handlers
	}

	// Parse file (simplified - in reality would parse Go code or use templates)
	// For now, return a default handler that returns the file content
	handlers["GET"] = func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Handler",
			"path":    c.Request.URL.Path,
			"file":    filePath,
		})
	}

	// Check if file has specific method handlers
	// In a full implementation, this would parse the Go file
	// and extract handler functions for each HTTP method

	return handlers
}

// fileToURL converts a file path to URL path
func fileToURL(filePath string) string {
	// Remove extension
	url := strings.TrimSuffix(filePath, filepath.Ext(filePath))

	// Convert to URL path
	url = filepath.ToSlash(url)

	// Handle dynamic routes [id] -> :id
	url = strings.ReplaceAll(url, "[", ":")
	url = strings.ReplaceAll(url, "]", "")

	// Handle catch-all routes [...slug]
	url = strings.ReplaceAll(url, "...", "*")

	// Add leading slash
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	return url
}

// isHandlerFile checks if a file is an API handler
func isHandlerFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".go" || ext == ".go.html" || ext == ".json"
}

// Helper functions

func sendJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func sendError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": message,
	})
}

func parseJSON(c *gin.Context, v interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

// fileExists checks if a file exists
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
