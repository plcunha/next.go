package static

import (
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// StaticConfig represents static file serving configuration
type StaticConfig struct {
	Dir       string
	Prefix    string
	MaxAge    time.Duration
	Immutable bool
}

// ServeStatic serves static files from a directory
func ServeStatic(config StaticConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if path starts with prefix
		if config.Prefix != "" && !strings.HasPrefix(c.Request.URL.Path, config.Prefix) {
			c.Next()
			return
		}

		// Remove prefix from path
		path := c.Request.URL.Path
		if config.Prefix != "" {
			path = strings.TrimPrefix(path, config.Prefix)
		}

		// Clean the path to prevent directory traversal
		path = filepath.Clean(path)
		fullPath := filepath.Join(config.Dir, path)

		// Check if file exists
		info, err := os.Stat(fullPath)
		if err != nil {
			c.Next()
			return
		}

		// Don't serve directories
		if info.IsDir() {
			c.Next()
			return
		}

		// Set cache headers
		if config.MaxAge > 0 {
			cacheControl := fmt.Sprintf("public, max-age=%.0f", config.MaxAge.Seconds())
			if config.Immutable {
				cacheControl += ", immutable"
			}
			c.Header("Cache-Control", cacheControl)
		}

		// Set content type
		ext := filepath.Ext(fullPath)
		contentType := mime.TypeByExtension(ext)
		if contentType != "" {
			c.Header("Content-Type", contentType)
		}

		// Serve the file
		c.File(fullPath)
		c.Abort()
	}
}

// ServePublic serves files from the public directory
func ServePublic(appDir string) gin.HandlerFunc {
	publicDir := filepath.Join(appDir, "public")
	
	return func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/") {
			c.Next()
			return
		}

		// Remove leading slash
		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		fullPath := filepath.Join(publicDir, path)

		// Check if file exists
		info, err := os.Stat(fullPath)
		if err != nil {
			c.Next()
			return
		}

		if info.IsDir() {
			c.Next()
			return
		}

		c.File(fullPath)
		c.Abort()
	}
}

// ServeNextStatic serves .next/static files
func ServeNextStatic(appDir string) gin.HandlerFunc {
	staticDir := filepath.Join(appDir, ".next", "static")
	
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		
		// Match /_next/static/*
		if !strings.HasPrefix(path, "/_next/static/") {
			c.Next()
			return
		}

		// Extract file path
		filePath := strings.TrimPrefix(path, "/_next/static/")
		fullPath := filepath.Join(staticDir, filePath)

		// Check if file exists
		if _, err := os.Stat(fullPath); err != nil {
			c.Next()
			return
		}

		// Set long cache for static assets
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		
		c.File(fullPath)
		c.Abort()
	}
}

// FileSystem is a wrapper for fs.FS to implement http.FileSystem
type FileSystem struct {
	fs fs.FS
}

// Open opens a file from the filesystem
func (f FileSystem) Open(name string) (http.File, error) {
	return f.fs.Open(name)
}

// WalkDir walks through a directory and calls fn for each file
func WalkDir(root string, fn func(path string, info os.FileInfo) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return fn(path, info)
	})
}

// CopyStaticFiles copies static files from src to dst
func CopyStaticFiles(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, 0644)
	})
}
