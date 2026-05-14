package apihandler

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// DevServer manages a compiled API server subprocess for dev mode
type DevServer struct {
	appDir    string
	tmpDir    string
	cmd       *exec.Cmd
	port      int
	proxy     *httputil.ReverseProxy
	mu        sync.Mutex
	handlers  []handlerInfo
	nextgoDir string
}

type handlerInfo struct {
	relPath string
	route   string
}

// New creates a new dev API server
func New(appDir string) *DevServer {
	return &DevServer{
		appDir:    appDir,
		tmpDir:    filepath.Join(appDir, ".next", "dev-api"),
		nextgoDir: nextgoSourceDir(),
	}
}

// Start compiles and starts the API server subprocess
func (d *DevServer) Start() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.scanHandlers(); err != nil {
		return fmt.Errorf("scan handlers: %w", err)
	}

	if len(d.handlers) == 0 {
		return nil // No API routes
	}

	// Find available port
	port, err := findFreePort()
	if err != nil {
		return err
	}
	d.port = port

	if err := d.generateAndCompile(port); err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	// Start subprocess
	d.cmd = exec.Command(filepath.Join(d.tmpDir, "api-server"))
	d.cmd.Dir = d.tmpDir
	d.cmd.Stdout = os.Stdout
	d.cmd.Stderr = os.Stderr
	if err := d.cmd.Start(); err != nil {
		return fmt.Errorf("start subprocess: %w", err)
	}

	// Setup proxy
	d.proxy = httputil.NewSingleHostReverseProxy(
		&url.URL{Scheme: "http", Host: fmt.Sprintf("127.0.0.1:%d", port)},
	)

	// Wait briefly for the subprocess to be ready
	time.Sleep(300 * time.Millisecond)

	return nil
}

// Stop stops the API server subprocess
func (d *DevServer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cmd != nil && d.cmd.Process != nil {
		d.cmd.Process.Kill()
		d.cmd = nil
	}
}

// Proxy returns the reverse proxy handler, or nil if no API server
func (d *DevServer) Proxy() http.Handler {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.proxy
}

// Reload recompiles the API server (called on file change)
func (d *DevServer) Reload() error {
	d.Stop()
	return d.Start()
}

// HasHandlers returns true if there are API routes
func (d *DevServer) HasHandlers() bool {
	return len(d.handlers) > 0
}

func (d *DevServer) scanHandlers() error {
	apiDir := filepath.Join(d.appDir, "app", "api")
	if _, err := os.Stat(apiDir); os.IsNotExist(err) {
		d.handlers = nil
		return nil
	}

	d.handlers = nil
	return filepath.Walk(apiDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".go.html") {
			return nil
		}
		relPath, _ := filepath.Rel(apiDir, path)
		relPath = filepath.ToSlash(relPath)
		dir := filepath.Dir(relPath)
		route := "/api"
		if dir != "." {
			route = "/api/" + filepath.ToSlash(dir)
		}
		d.handlers = append(d.handlers, handlerInfo{relPath: relPath, route: route})
		return nil
	})
}

func (d *DevServer) generateAndCompile(port int) error {
	os.MkdirAll(d.tmpDir, 0755)

	// Copy handler files, changing package main to correct package name
	for _, h := range d.handlers {
		src := filepath.Join(d.appDir, "app", "api", filepath.FromSlash(h.relPath))
		dst := filepath.Join(d.tmpDir, "handlers", h.relPath)
		os.MkdirAll(filepath.Dir(dst), 0755)
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		dir := filepath.Dir(h.relPath)
		pkgName := filepath.Base(dir)
		if pkgName == "handlers" || pkgName == "." || pkgName == "" {
			pkgName = "root"
		}
		pkgName = sanitizePkgName(pkgName)
		content := strings.Replace(string(data), "package main", "package "+pkgName, 1)
		os.WriteFile(dst, []byte(content), 0644)
	}

	// Generate main.go
	mainContent := d.generateMain(port)
	if err := os.WriteFile(filepath.Join(d.tmpDir, "main.go"), []byte(mainContent), 0644); err != nil {
		return err
	}

	// Generate go.mod with proper replace directive
	nextgoPath := d.nextgoDir
	if nextgoPath == "" {
		nextgoPath = detectNextGoPathFallback()
	}
	if nextgoPath == "" {
		return fmt.Errorf("cannot find next.go source — set NEXTGO_SRC env var")
	}
	absPath, _ := filepath.Abs(nextgoPath)

	goMod := fmt.Sprintf(`module api-server-dev

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/nextgo/nextgo v0.0.0
)

replace github.com/nextgo/nextgo => %s
`, absPath)
	if err := os.WriteFile(filepath.Join(d.tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return err
	}

	// Copy parent go.sum
	parentGoSum := filepath.Join(d.nextgoDir, "go.sum")
	if data, err := os.ReadFile(parentGoSum); err == nil {
		os.WriteFile(filepath.Join(d.tmpDir, "go.sum"), data, 0644)
	}

	// go mod tidy (-e ignores packages that can't load, e.g. test deps
	// needing newer Go toolchains)
	tidy := exec.Command("go", "mod", "tidy", "-e")
	tidy.Dir = d.tmpDir
	tidy.Stdout = os.Stdout
	tidy.Stderr = os.Stderr
	tidy.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	if err := tidy.Run(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	// Build
	cmd := exec.Command("go", "build", "-o", "api-server", ".")
	cmd.Dir = d.tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return nil
}

func (d *DevServer) generateMain(port int) string {
	var sb strings.Builder
	sb.WriteString("package main\n\nimport (\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"os\"\n")
	sb.WriteString("\t\"os/signal\"\n")
	sb.WriteString("\t\"syscall\"\n\n")
	sb.WriteString("\t\"github.com/gin-gonic/gin\"\n")

	// Handler imports
	for _, h := range d.handlers {
		dir := filepath.Dir(h.relPath)
		pkgName := filepath.Base(dir)
		if pkgName == "handlers" || pkgName == "." || pkgName == "" {
			pkgName = "root"
		}
		pkgName = sanitizePkgName(pkgName)
		importPath := "api-server-dev/handlers/" + dir
		if importPath == "api-server-dev/handlers/" {
			importPath = "api-server-dev/handlers"
		}
		sb.WriteString(fmt.Sprintf("\t%s \"%s\"\n", pkgName, importPath))
	}

	sb.WriteString(")\n\n")

	sb.WriteString("func main() {\n")
	sb.WriteString("\tgin.SetMode(gin.ReleaseMode)\n")
	sb.WriteString("\tr := gin.New()\n")
	sb.WriteString("\tr.Use(gin.Recovery())\n\n")

	// Register routes
	for _, h := range d.handlers {
		dir := filepath.Dir(h.relPath)
		pkgName := filepath.Base(dir)
		if pkgName == "handlers" || pkgName == "." || pkgName == "" {
			pkgName = "root"
		}
		pkgName = sanitizePkgName(pkgName)
		sb.WriteString(fmt.Sprintf("\tr.Any(\"%s\", %s.Handler)\n", h.route, pkgName))
	}

	sb.WriteString(fmt.Sprintf("\n\tport := %d\n", port))
	sb.WriteString("\tfmt.Printf(\"API server ready on :%d\\n\", port)\n\n")
	sb.WriteString("\t// Graceful shutdown\n")
	sb.WriteString("\tsig := make(chan os.Signal, 1)\n")
	sb.WriteString("\tsignal.Notify(sig, os.Interrupt, syscall.SIGTERM)\n")
	sb.WriteString("\tgo func() {\n")
	sb.WriteString("\t\t<-sig\n")
	sb.WriteString("\t\tos.Exit(0)\n")
	sb.WriteString("\t}()\n\n")
	sb.WriteString("\tr.Run(fmt.Sprintf(\":%d\", port))\n")
	sb.WriteString("}\n")

	return sb.String()
}

func findFreePort() (int, error) {
	for port := 4600; port < 4700; port++ {
		if isPortFree(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no free port found in range 4600-4699")
}

func isPortFree(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 100*time.Millisecond)
	if err != nil {
		return true
	}
	conn.Close()
	return false
}

// nextgoSourceDir finds the next.go module root using runtime.Caller.
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

func detectNextGoPathFallback() string {
	if env := os.Getenv("NEXTGO_SRC"); env != "" {
		if _, err := os.Stat(filepath.Join(env, "go.mod")); err == nil {
			return env
		}
	}
	return ""
}

func sanitizePkgName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	if len(name) == 0 || (name[0] >= '0' && name[0] <= '9') {
		name = "pkg" + name
	}
	return name
}
