package router

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// RouteType represents the type of route
type RouteType string

const (
	RouteTypePage    RouteType = "page"
	RouteTypeLayout  RouteType = "layout"
	RouteTypeLoading RouteType = "loading"
	RouteTypeError   RouteType = "error"
	RouteTypeAPI     RouteType = "api"
	RouteTypeNotFound RouteType = "not-found"
)

// Route represents a file-system route
type Route struct {
	Path     string
	FilePath string
	Type     RouteType
	Methods  []string
	Metadata map[string]interface{}
}

// Router handles file-system based routing
type Router struct {
	AppDir  string
	PagesDir string
	Routes  map[string]*Route
	mu      sync.RWMutex
}

// New creates a new router
func New(appDir string) *Router {
	pagesDir := filepath.Join(appDir, "app")
	if _, err := os.Stat(pagesDir); os.IsNotExist(err) {
		pagesDir = appDir
	}

	return &Router{
		AppDir:   appDir,
		PagesDir: pagesDir,
		Routes:   make(map[string]*Route),
	}
}

// Scan scans the app directory and builds the route tree
func (r *Router) Scan() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Routes = make(map[string]*Route)

	err := filepath.Walk(r.PagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and files
		if strings.HasPrefix(info.Name(), ".") || strings.HasPrefix(info.Name(), "_") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip node_modules, .next, etc.
		if info.Name() == "node_modules" || info.Name() == ".next" || info.Name() == "components" || info.Name() == "lib" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			r.processFile(path)
		}

		return nil
	})

	return err
}

// processFile processes a file and creates routes
func (r *Router) processFile(path string) {
	relPath, _ := filepath.Rel(r.PagesDir, path)
	dir := filepath.Dir(relPath)
	filename := filepath.Base(path)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	ext := filepath.Ext(filename)

	// Skip non-relevant files
	if ext != ".go.html" && ext != ".html" && ext != ".go" {
		return
	}

	// Skip non-route files (layout, loading, error are NOT routes in App Router)
	if name == "layout" || name == "loading" || name == "error" || name == "not-found" {
		return
	}

	// Determine route type
	var routeType RouteType
	switch {
	case name == "page":
		routeType = RouteTypePage
	case strings.Contains(relPath, string(os.PathSeparator)+"api"+string(os.PathSeparator)):
		routeType = RouteTypeAPI
	default:
		// Custom page (e.g., about.go.html creates /about)
		routeType = RouteTypePage
	}

	// Calculate URL path from directory structure
	urlPath := r.dirToURLPath(dir)

	// Handle dynamic routes [id] or [...slug]
	urlPath = r.processDynamicSegments(urlPath)

	// Only create route if it's a page or API route
	if routeType != RouteTypePage && routeType != RouteTypeAPI {
		return
	}

	// Get or create route
	route, exists := r.Routes[urlPath]
	if !exists {
		route = &Route{
			Path:     urlPath,
			Type:     routeType,
			Methods:  []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
			Metadata: make(map[string]interface{}),
		}
		r.Routes[urlPath] = route
	}

	route.FilePath = path
	route.Type = routeType

	fmt.Printf("  ✓ Route: %s -> %s (%s)\n", urlPath, path, routeType)
}

// dirToURLPath converts a directory path to URL path
func (r *Router) dirToURLPath(dir string) string {
	if dir == "." {
		return "/"
	}

	// Convert to URL path
	urlPath := strings.ReplaceAll(dir, string(os.PathSeparator), "/")
	urlPath = strings.Trim(urlPath, "/")

	if urlPath == "" {
		return "/"
	}

	return "/" + urlPath
}

// processDynamicSegments processes dynamic route segments
func (r *Router) processDynamicSegments(path string) string {
	// Handle catch-all [...slug] → *slug (must come before [id] → :id)
	path = strings.ReplaceAll(path, "[...", "*")
	// Replace [id] with :id
	path = strings.ReplaceAll(path, "[", ":")
	path = strings.ReplaceAll(path, "]", "")
	return path
}

// GetRoute returns a route by path
func (r *Router) GetRoute(path string) (*Route, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	route, exists := r.Routes[path]
	return route, exists
}

// GetAllRoutes returns all routes
func (r *Router) GetAllRoutes() map[string]*Route {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routes := make(map[string]*Route)
	for k, v := range r.Routes {
		routes[k] = v
	}
	return routes
}

// Close cleans up resources
func (r *Router) Close() {
	// Cleanup if needed
}
