package build

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

// Builder handles building the Next.go application
type Builder struct {
	AppDir    string
	OutDir    string
	StartAt   time.Time
	pages     []string
	files     map[string]string
	nextgoDir string // resolved source directory of next.go itself
}

// New creates a new builder
func New(appDir string) *Builder {
	return &Builder{
		AppDir:    appDir,
		OutDir:    filepath.Join(appDir, ".next"),
		pages:     make([]string, 0),
		files:     make(map[string]string),
		nextgoDir: nextgoSourceDir(),
	}
}

// nextgoSourceDir finds the next.go module root using runtime.Caller.
// This is the reliable way to locate the source regardless of where the
// binary is installed — no GOPATH or hardcoded paths needed.
func nextgoSourceDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			abs, _ := filepath.Abs(dir)
			return abs
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// Build builds the application for production
func (b *Builder) Build() error {
	b.StartAt = time.Now()

	fmt.Println("Building for production...")

	if err := b.ensureDirs(); err != nil {
		return err
	}

	if err := b.buildPages(); err != nil {
		return err
	}

	if err := b.buildAPIRoutes(); err != nil {
		return err
	}

	if err := b.buildStaticAssets(); err != nil {
		return err
	}

	if err := b.generateManifest(); err != nil {
		return err
	}

	buildID := fmt.Sprintf("%d", time.Now().Unix())

	fmt.Printf("\n✓ Build complete in %v\n", time.Since(b.StartAt))
	fmt.Printf("✓ Build ID: %s\n", buildID)
	fmt.Printf("✓ Output directory: %s\n", b.OutDir)

	return nil
}

// ensureDirs ensures all necessary directories exist
func (b *Builder) ensureDirs() error {
	dirs := []string{
		"server",
		"static",
		"server/app",
		"static/chunks",
		"cache",
		"server/pages",
		"build",
		"build/handlers",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(b.OutDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

// buildPages builds all pages
func (b *Builder) buildPages() error {
	fmt.Println("  Building pages...")

	return filepath.Walk(filepath.Join(b.AppDir, "app"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".go.html") && !strings.HasSuffix(path, ".html") {
			return nil
		}

		filename := filepath.Base(path)
		name := strings.TrimSuffix(filename, filepath.Ext(filename))

		if name == "layout" || name == "loading" || name == "error" {
			return nil
		}

		return b.buildPage(path)
	})
}

// buildPage builds a single page
func (b *Builder) buildPage(pagePath string) error {
	relPath, _ := filepath.Rel(filepath.Join(b.AppDir, "app"), pagePath)
	pageName := filepath.ToSlash(strings.TrimSuffix(relPath, filepath.Ext(relPath)))
	pageName = strings.TrimSuffix(pageName, ".go")

	content, err := os.ReadFile(pagePath)
	if err != nil {
		return err
	}

	// Find and apply layout
	var finalContent []byte
	layoutContent := b.findLayoutForPage(pagePath)
	if layoutContent != nil {
		combined := strings.Replace(string(layoutContent), "{{.children}}", string(content), 1)
		finalContent = []byte(combined)
	} else {
		finalContent = content
	}

	// Minify
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)

	minified, err := m.String("text/html", string(finalContent))
	if err != nil {
		minified = string(finalContent)
	}

	outputPath := filepath.Join(b.OutDir, "server", "app", pageName+".html")
	os.MkdirAll(filepath.Dir(outputPath), 0755)

	if err := os.WriteFile(outputPath, []byte(minified), 0644); err != nil {
		return err
	}

	urlPath := "/" + strings.TrimSuffix(pageName, "/page")
	if urlPath == "//" {
		urlPath = "/"
	}
	b.pages = append(b.pages, urlPath)
	b.files[pageName+".html"] = "server/app/" + pageName + ".html"

	fmt.Printf("    ✓ %s\n", pageName)
	return nil
}

// findLayoutForPage walks up from the page to find a layout.go.html
func (b *Builder) findLayoutForPage(pagePath string) []byte {
	dir := filepath.Dir(pagePath)
	appDir := filepath.Join(b.AppDir, "app")

	for {
		layoutPath := filepath.Join(dir, "layout.go.html")
		if content, err := os.ReadFile(layoutPath); err == nil {
			return content
		}
		parent := filepath.Dir(dir)
		if parent == dir || parent == appDir || !strings.HasPrefix(parent, appDir) {
			break
		}
		dir = parent
	}
	return nil
}

// handlerInfo describes a discovered API handler
type handlerInfo struct {
	relPath string
	route   string
	pkgName string
	dir     string
}

// buildAPIRoutes scans API handlers, generates Go code and compiles
func (b *Builder) buildAPIRoutes() error {
	fmt.Println("  Building API routes...")

	apiDir := filepath.Join(b.AppDir, "app", "api")
	if _, err := os.Stat(apiDir); os.IsNotExist(err) {
		fmt.Println("    (no API routes)")
		return nil
	}

	var handlers []handlerInfo

	err := filepath.Walk(apiDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".go.html") {
			return nil
		}
		relPath, _ := filepath.Rel(apiDir, path)
		relPath = filepath.ToSlash(relPath)
		hDir := filepath.Dir(relPath)
		route := "/api"
		if hDir != "." {
			route = "/api/" + filepath.ToSlash(hDir)
		}
		pkgName := filepath.Base(hDir)
		if pkgName == "handlers" || pkgName == "." || pkgName == "" {
			pkgName = "root"
		}
		handlers = append(handlers, handlerInfo{
			relPath: relPath,
			route:   route,
			pkgName: sanitizePkgName(pkgName),
			dir:     hDir,
		})
		fmt.Printf("    ✓ API: %s → %s\n", relPath, route)
		b.files["api/"+relPath] = route
		return nil
	})
	if err != nil {
		return err
	}

	if len(handlers) == 0 {
		return nil
	}

	// Copy handler files, changing package main to correct package name
	for _, h := range handlers {
		src := filepath.Join(apiDir, filepath.FromSlash(h.relPath))
		dst := filepath.Join(b.OutDir, "build", "handlers", h.dir, "handler.go")
		os.MkdirAll(filepath.Dir(dst), 0755)
		data, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("read handler %s: %w", h.relPath, err)
		}
		content := strings.Replace(string(data), "package main", "package "+h.pkgName, 1)
		if err := os.WriteFile(dst, []byte(content), 0644); err != nil {
			return err
		}
	}

	// Generate main.go for the build project
	if err := b.generateBuildMain(handlers); err != nil {
		return err
	}

	// Generate go.mod with proper replace to the detected next.go source
	if err := b.generateBuildGoMod(); err != nil {
		return err
	}

	// Copy parent go.sum as a starting point so go mod tidy doesn't
	// have to resolve everything from scratch
	parentGoSum := filepath.Join(b.nextgoDir, "go.sum")
	if data, err := os.ReadFile(parentGoSum); err == nil {
		os.WriteFile(filepath.Join(b.OutDir, "build", "go.sum"), data, 0644)
	}

	// Tidy modules and download dependencies.
	// The -e flag ignores packages that can't be loaded (e.g., test deps
	// that require newer Go toolchains), which keeps builds working on
	// machines without the latest Go.
	fmt.Println("  Resolving dependencies...")
	tidyCmd := exec.Command("go", "mod", "tidy", "-e")
	tidyCmd.Dir = filepath.Join(b.OutDir, "build")
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	tidyCmd.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	if err := tidyCmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	// Compile
	fmt.Println("  Compiling production server...")
	cmd := exec.Command("go", "build", "-o", filepath.Join(b.OutDir, "server", "app.exe"), ".")
	cmd.Dir = filepath.Join(b.OutDir, "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("compile production server: %w (Go must be installed and on PATH)", err)
	}

	fmt.Println("    ✓ Compiled app.exe")
	return nil
}

// generateBuildMain creates main.go for the compiled production server.
// It imports all handler packages and registers them, serves pre-rendered
// pages from .next/server/app/, and serves static assets.
func (b *Builder) generateBuildMain(handlers []handlerInfo) error {
	var sb strings.Builder
	sb.WriteString("package main\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"net/http\"\n")
	sb.WriteString("\t\"os\"\n")
	sb.WriteString("\t\"path/filepath\"\n")
	sb.WriteString("\t\"strings\"\n")
	sb.WriteString("\n")
	sb.WriteString("\t\"github.com/gin-gonic/gin\"\n")
	sb.WriteString("\t\"github.com/nextgo/nextgo/packages/middleware\"\n")

	// Import handler packages
	for _, h := range handlers {
		importPath := "nextgo-build/handlers/" + h.dir
		if importPath == "nextgo-build/handlers/" {
			importPath = "nextgo-build/handlers"
		}
		sb.WriteString(fmt.Sprintf("\t%s \"%s\"\n", h.pkgName, importPath))
	}
	sb.WriteString(")\n\n")

	sb.WriteString("func main() {\n")
	sb.WriteString("\tgin.SetMode(gin.ReleaseMode)\n")
	sb.WriteString("\tr := gin.New()\n\n")

	// Middleware
	sb.WriteString("\t// Middleware chain\n")
	sb.WriteString("\tr.Use(gin.Recovery())\n")
	sb.WriteString("\tr.Use(middleware.Logger())\n")
	sb.WriteString("\tr.Use(middleware.CORS())\n")
	sb.WriteString("\tr.Use(middleware.Security())\n\n")

	// Static assets
	sb.WriteString("\t// Static assets\n")
	sb.WriteString("\tbuildDir := filepath.Dir(os.Args[0])\n")
	sb.WriteString("\tstaticDir := filepath.Join(buildDir, \"..\", \"static\")\n")
	sb.WriteString("\tr.StaticFS(\"/_next/static\", http.Dir(staticDir))\n")
	sb.WriteString("\tr.StaticFS(\"/public\", http.Dir(staticDir))\n\n")

	// Serve pre-rendered pages (SSG)
	sb.WriteString("\t// Pre-rendered pages (SSG)\n")
	sb.WriteString("\tpagesDir := filepath.Join(buildDir, \"..\", \"server\", \"app\")\n")
	sb.WriteString("\tr.NoRoute(func(c *gin.Context) {\n")
	sb.WriteString("\t\tpath := c.Request.URL.Path\n")
	sb.WriteString("\t\tif path == \"/\" {\n")
	sb.WriteString("\t\t\tpath = \"/page\"\n")
	sb.WriteString("\t\t}\n")
	sb.WriteString("\t\tclean := strings.Trim(path, \"/\")\n")
	sb.WriteString("\t\tfilePath := filepath.Join(pagesDir, clean+\".html\")\n")
	sb.WriteString("\t\tif _, err := os.Stat(filePath); err != nil {\n")
	sb.WriteString("\t\t\tfilePath = filepath.Join(pagesDir, clean, \"page.html\")\n")
	sb.WriteString("\t\t}\n")
	sb.WriteString("\t\tdata, err := os.ReadFile(filePath)\n")
	sb.WriteString("\t\tif err != nil {\n")
	sb.WriteString("\t\t\tc.String(http.StatusNotFound, \"404 - Page not found\")\n")
	sb.WriteString("\t\t\treturn\n")
	sb.WriteString("\t\t}\n")
	sb.WriteString("\t\tc.Data(http.StatusOK, \"text/html; charset=utf-8\", data)\n")
	sb.WriteString("\t})\n\n")

	// Register API handlers
	for _, h := range handlers {
		sb.WriteString(fmt.Sprintf("\t// Route: %s\n", h.route))
		sb.WriteString(fmt.Sprintf("\tr.Any(\"%s\", %s.Handler)\n\n", h.route, h.pkgName))
	}

	sb.WriteString("\t// Start server\n")
	sb.WriteString("\tport := os.Getenv(\"PORT\")\n")
	sb.WriteString("\tif port == \"\" {\n")
	sb.WriteString("\t\tport = \"3000\"\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\tfmt.Printf(\"🚀 Next.go production server on :%s\\n\", port)\n")
	sb.WriteString("\tr.Run(\":\" + port)\n")
	sb.WriteString("}\n")

	mainPath := filepath.Join(b.OutDir, "build", "main.go")
	return os.WriteFile(mainPath, []byte(sb.String()), 0644)
}

// generateBuildGoMod creates go.mod for the build project.
// Uses the detected next.go source directory for the replace directive.
func (b *Builder) generateBuildGoMod() error {
	nextgoPath := b.nextgoDir
	if nextgoPath == "" {
		// Fallback: try common locations
		nextgoPath = detectNextGoPathFallback()
	}
	if nextgoPath == "" {
		return fmt.Errorf("cannot find next.go source directory — run from within the next.go repository or set NEXTGO_SRC")
	}

	// Use absolute path for the replace directive
	absPath, err := filepath.Abs(nextgoPath)
	if err != nil {
		return fmt.Errorf("resolve next.go path: %w", err)
	}

	goMod := fmt.Sprintf(`module nextgo-build

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/nextgo/nextgo v0.0.0
)

replace github.com/nextgo/nextgo => %s
`, absPath)

	goModPath := filepath.Join(b.OutDir, "build", "go.mod")
	return os.WriteFile(goModPath, []byte(goMod), 0644)
}

// detectNextGoPathFallback tries common locations when runtime.Caller fails.
func detectNextGoPathFallback() string {
	// Check environment variable first
	if env := os.Getenv("NEXTGO_SRC"); env != "" {
		if _, err := os.Stat(filepath.Join(env, "go.mod")); err == nil {
			return env
		}
	}

	// Check relative to the project being built (common case: next.go source is nearby)
	candidates := []string{
		"../../..",
		"../..",
		"..",
	}

	for _, c := range candidates {
		abs, _ := filepath.Abs(c)
		if _, err := os.Stat(filepath.Join(abs, "go.mod")); err == nil {
			// Verify it's the next.go module
			data, err := os.ReadFile(filepath.Join(abs, "go.mod"))
			if err == nil && strings.Contains(string(data), "github.com/nextgo/nextgo") {
				return abs
			}
		}
	}

	return ""
}

// buildStaticAssets builds and optimizes static assets
func (b *Builder) buildStaticAssets() error {
	fmt.Println("  Building static assets...")

	publicDir := filepath.Join(b.AppDir, "public")
	if _, err := os.Stat(publicDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(publicDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(publicDir, path)
		outputPath := filepath.Join(b.OutDir, "static", relPath)

		if info.IsDir() {
			return os.MkdirAll(outputPath, 0755)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		os.MkdirAll(filepath.Dir(outputPath), 0755)
		b.files["static/"+relPath] = "static/" + relPath
		return os.WriteFile(outputPath, data, 0644)
	})
}

// generateManifest generates the build manifest
func (b *Builder) generateManifest() error {
	fmt.Println("  Generating build manifest...")

	manifest := map[string]interface{}{
		"version": "0.2.0",
		"buildId": fmt.Sprintf("%d", time.Now().Unix()),
		"pages":   b.pages,
		"files":   b.files,
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	manifestPath := filepath.Join(b.OutDir, "build-manifest.json")
	return os.WriteFile(manifestPath, data, 0644)
}

// sanitizePkgName converts a directory name to a valid Go package name.
func sanitizePkgName(name string) string {
	// Replace hyphens with underscores, remove other invalid chars
	name = strings.ReplaceAll(name, "-", "_")
	// Ensure it's a valid Go identifier
	if len(name) == 0 || (name[0] >= '0' && name[0] <= '9') {
		name = "pkg" + name
	}
	return name
}
