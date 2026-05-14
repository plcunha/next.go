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
	"github.com/nextgo/nextgo/packages/apihandler"
	"github.com/nextgo/nextgo/packages/config"
	"github.com/nextgo/nextgo/packages/middleware"
	"github.com/nextgo/nextgo/packages/router"
	"github.com/nextgo/nextgo/packages/watcher"
)

// PagePropsFunc is a function that fetches data for a page before rendering.
// Pages can optionally implement this pattern by placing a .go file alongside
// their .go.html template that exports a Props function.
//
// Example (app/blog/[id]/page.go alongside page.go.html):
//
//	package blog
//	func Props(c *gin.Context) gin.H {
//	    id := c.Param("id")
//	    post := fetchPost(id)
//	    return gin.H{"post": post, "title": post.Title}
//	}
type PagePropsFunc func(c *gin.Context) gin.H

// Server represents the Next.go server
type Server struct {
	Router     *router.Router
	AppDir     string
	Port       string
	DevMode    bool
	gin        *gin.Engine
	apiDev     *apihandler.DevServer
	mu         sync.RWMutex
	propsFuncs map[string]PagePropsFunc
}

// New creates a new server instance
func New(appDir string) *Server {
	s := &Server{
		AppDir:     appDir,
		Port:       ":3000",
		DevMode:    true,
		propsFuncs: make(map[string]PagePropsFunc),
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

	s.gin = gin.New()

	// Middleware chain
	s.gin.Use(gin.Recovery())
	if s.DevMode {
		s.gin.Use(s.loggingMiddleware())
	} else {
		s.gin.Use(middleware.Logger())
	}
	s.gin.Use(middleware.CORS())
	s.gin.Use(middleware.Security())

	// Start API dev server (compiles handlers into subprocess)
	if s.DevMode {
		s.apiDev = apihandler.New(s.AppDir)
		if err := s.apiDev.Start(); err != nil {
			fmt.Printf("  ⚠️ API dev server: %v (handlers will return placeholder)\n", err)
		} else if s.apiDev.HasHandlers() {
			fmt.Printf("  ✓ API handlers compiled and running\n")
		}
	}

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
			if s.apiDev != nil {
				if err := s.apiDev.Reload(); err != nil {
					fmt.Printf("  ⚠️ API reload: %v\n", err)
				}
			}
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

// renderPage renders a page with SSR, including optional data fetching
// via getServerSideProps pattern (Props() function in page package).
func (s *Server) renderPage(c *gin.Context, route *router.Route) {
	content, err := os.ReadFile(route.FilePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading file: %v", err)
		return
	}

	pageContent := string(content)

	// Apply layout
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

	// Prepare data — base template variables
	data := gin.H{
		"title": titleFromPath(route.Path),
		"path":  c.Request.URL.Path,
		"route": route.Path,
	}

	// Check for getServerSideProps — look for a .go file alongside the page
	// that exports a Props function. In a full implementation this would
	// use Go plugins or code generation; for now we enable the pattern via
	// explicit registration.
	if propsFunc, ok := s.propsFuncs[route.Path]; ok {
		props := propsFunc(c)
		for k, v := range props {
			data[k] = v
		}
	}

	// Execute template
	w := c.Writer
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = tmpl.Execute(w, data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Template execute error: %v", err)
	}
}

// RegisterProps registers a getServerSideProps function for a route.
// This allows pages to fetch data at request time before rendering.
func (s *Server) RegisterProps(routePath string, fn PagePropsFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.propsFuncs[routePath] = fn
}

// renderAPI renders an API route
func (s *Server) renderAPI(c *gin.Context, route *router.Route) {
	if s.apiDev != nil && s.apiDev.HasHandlers() {
		if proxy := s.apiDev.Proxy(); proxy != nil {
			proxy.ServeHTTP(c.Writer, c.Request)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"route":   route.Path,
		"method":  c.Request.Method,
		"handler": route.FilePath,
	})
}

// findLayout finds the nearest layout file
func (s *Server) findLayout(dir string) string {
	for {
		layoutPath := filepath.Join(dir, "layout.go.html")
		if _, err := os.Stat(layoutPath); err == nil {
			content, _ := os.ReadFile(layoutPath)
			return string(content)
		}

		parent := filepath.Dir(dir)
		if parent == dir || parent == s.Router.PagesDir || parent == "." {
			break
		}
		dir = parent
	}

	return `<!DOCTYPE html>
<html><head><title>{{.title}}</title></head>
<body>{{.children}}</body></html>`
}

// handleRequest handles requests that don't match registered routes
func (s *Server) handleRequest(c *gin.Context) {
	path := c.Request.URL.Path

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

// titleFromPath derives a page title from a route path
func titleFromPath(path string) string {
	if path == "/" {
		return "Home"
	}
	p := strings.Trim(path, "/")
	if last := strings.LastIndex(p, "/"); last >= 0 {
		p = p[last+1:]
	}
	p = strings.ReplaceAll(p, "-", " ")
	p = strings.ReplaceAll(p, "_", " ")
	return strings.Title(p)
}
