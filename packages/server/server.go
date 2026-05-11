package server

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nextgo/nextgo/packages/config"
	"github.com/nextgo/nextgo/packages/router"
	"github.com/nextgo/nextgo/packages/watcher"
)

// Server represents the Next.go server
type Server struct {
	Router  *router.Router
	AppDir  string
	Port    string
	DevMode bool
	gin     *gin.Engine
	mu      sync.RWMutex
}

// New creates a new server instance
func New(appDir string) *Server {
	s := &Server{
		AppDir:  appDir,
		Port:    ":3000",
		DevMode: true,
	}

	// Load config if exists
	if cfg, err := config.Load(appDir); err == nil {
		s.Port = fmt.Sprintf(":%d", cfg.Server.Port)
	}

	// Initialize router
	s.Router = router.New(appDir)
	s.Router.Scan()

	return s
}

// Start starts the server
func (s *Server) Start(port string) error {
	if port != "" {
		s.Port = port
	}

	if !s.DevMode {
		gin.SetMode(gin.ReleaseMode)
	}

	s.gin = gin.Default()

	// Add middleware
	s.gin.Use(s.loggingMiddleware())

	// Register file-system routes
	s.registerRoutes()

	// Static files
	s.gin.Static("/public", filepath.Join(s.AppDir, "public"))
	s.gin.Static("/_next/static", filepath.Join(s.AppDir, ".next", "static"))

	// Catch-all for SSR
	s.gin.NoRoute(s.handleRequest)

	// Start file watcher in dev mode
	if s.DevMode {
		go watcher.Watch(s.AppDir, func() {
			fmt.Println("🔄 File changed, re-scanning routes...")
			s.Router.Scan()
		})
	}

	fmt.Printf("\n🚀 Next.go server ready!\n")
	fmt.Printf("   Local:   http://localhost%s\n", s.Port)
	fmt.Printf("   App:     %s\n", s.Router.PagesDir)
	fmt.Printf("   Mode:    %s\n\n", map[bool]string{true: "development", false: "production"}[s.DevMode])

	return s.gin.Run(s.Port)
}

// registerRoutes registers all routes from the file system
func (s *Server) registerRoutes() {
	routes := s.Router.GetAllRoutes()

	for path, route := range routes {
		handler := s.createHandler(route)

		// Register for all methods
		for _, method := range route.Methods {
			switch method {
			case "GET":
				s.gin.GET(path, handler)
			case "POST":
				s.gin.POST(path, handler)
			case "PUT":
				s.gin.PUT(path, handler)
			case "DELETE":
				s.gin.DELETE(path, handler)
			case "PATCH":
				s.gin.PATCH(path, handler)
			}
		}
	}
}

// createHandler creates a gin handler for a route
func (s *Server) createHandler(route *router.Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		s.renderRoute(c, route)
	}
}

// renderRoute renders a route
func (s *Server) renderRoute(c *gin.Context, route *router.Route) {
	switch route.Type {
	case router.RouteTypePage:
		s.renderPage(c, route)
	case router.RouteTypeAPI:
		s.renderAPI(c, route)
	default:
		s.renderPage(c, route)
	}
}

// renderPage renders a page with SSR
func (s *Server) renderPage(c *gin.Context, route *router.Route) {
	// Read page content
	content, err := os.ReadFile(route.FilePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading file: %v", err)
		return
	}

	pageContent := string(content)

	// Try to find and apply layout
	layout := s.findLayout(filepath.Dir(route.FilePath))
	if layout != "" {
		pageContent = strings.Replace(layout, "{{.children}}", pageContent, 1)
	}

	// Parse and execute template
	tmpl, err := template.New("page").Parse(pageContent)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template parse error: %v", err)
		return
	}

	// Prepare data
	data := gin.H{
		"title": filepath.Base(route.Path),
		"path":  c.Request.URL.Path,
	}

	// Execute template
	w := c.Writer
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = tmpl.Execute(w, data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execute error: %v", err)
	}
}

// renderAPI renders an API route
func (s *Server) renderAPI(c *gin.Context, route *router.Route) {
	// For Go handler files, we would load and execute them
	// For now, return JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "API Route",
		"path":    route.Path,
		"method":  c.Request.Method,
	})
}

// findLayout finds the nearest layout file
func (s *Server) findLayout(dir string) string {
	// Walk up the directory tree to find layout.go.html
	for {
		layoutPath := filepath.Join(dir, "layout.go.html")
		if _, err := os.Stat(layoutPath); err == nil {
			content, _ := os.ReadFile(layoutPath)
			return string(content)
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir || parent == s.Router.PagesDir || parent == "." {
			break
		}
		dir = parent
	}

	// Return default layout
	return `<!DOCTYPE html>
<html><head><title>{{.title}}</title></head>
<body>{{.children}}</body></html>`
}

// handleRequest handles requests that don't match registered routes
func (s *Server) handleRequest(c *gin.Context) {
	path := c.Request.URL.Path

	// Try to find a matching page
	pagePath := filepath.Join(s.Router.PagesDir, path, "page.go.html")
	if _, err := os.Stat(pagePath); err == nil {
		route := &router.Route{
			Path:     path,
			FilePath: pagePath,
			Type:     router.RouteTypePage,
		}
		s.renderPage(c, route)
		return
	}

	// Try just the path as a file
	pagePath = filepath.Join(s.Router.PagesDir, path+".go.html")
	if _, err := os.Stat(pagePath); err == nil {
		route := &router.Route{
			Path:     path,
			FilePath: pagePath,
			Type:     router.RouteTypePage,
		}
		s.renderPage(c, route)
		return
	}

	// 404
	c.String(http.StatusNotFound, "404 - Page not found: %s", path)
}

// loggingMiddleware logs requests
func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		if s.DevMode {
			fmt.Printf("[%s] %s %d %v\n", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
		}
	}
}
