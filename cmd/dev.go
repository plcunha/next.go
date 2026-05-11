package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/nextgo/nextgo/packages/server"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start the development server with hot reload",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		port, _ := cmd.Flags().GetString("port")

		// Check if package.json exists, if not create a default Next.js-like project
		if !fileExists(filepath.Join(dir, "package.json")) {
			fmt.Println("Initializing Next.go project...")
			initProject(dir)
		}

		fmt.Println("Starting Next.go dev server...")

		// Start server
		s := server.New(dir)
		s.DevMode = true

		if port == "" {
			port = ":3000"
		} else if !strings.HasPrefix(port, ":") {
			port = ":" + port
		}

		// Open browser in dev mode
		go func() {
			time.Sleep(500 * time.Millisecond)
			openBrowser("http://localhost" + port)
		}()

		if err := s.Start(port); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	devCmd.Flags().StringP("port", "p", "", "Port to run the development server")
	rootCmd.AddCommand(devCmd)
}

func initProject(dir string) {
	// Create default project structure
	dirs := []string{
		"app",
		"app/api/hello",
		"app/about",
		"public",
		"components",
		"lib",
	}

	for _, d := range dirs {
		os.MkdirAll(filepath.Join(dir, d), 0755)
	}

	// Copy template files from template/default if exists, otherwise create basic ones
	createDefaultFiles(dir)

	fmt.Println("✓ Created Next.go project structure!")
}

func createDefaultFiles(dir string) {
	// Create layout
	layout := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.title}}</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 1200px; margin: 0 auto; padding: 2rem; }
        header { background: #000; color: white; padding: 1rem 0; }
        nav { max-width: 1200px; margin: 0 auto; padding: 0 2rem; display: flex; gap: 2rem; }
        nav a { color: white; text-decoration: none; }
        nav a:hover { text-decoration: underline; }
        main { min-height: calc(100vh - 200px); padding: 2rem 0; }
        footer { background: #f5f5f5; padding: 2rem 0; text-align: center; }
    </style>
</head>
<body>
    <header>
        <nav>
            <a href="/">Home</a>
            <a href="/about">About</a>
        </nav>
    </header>
    <main>
        <div class="container">{{.children}}</div>
    </main>
    <footer>
        <p>Built with Next.go</p>
    </footer>
</body>
</html>`

	os.WriteFile(filepath.Join(dir, "app", "layout.go.html"), []byte(layout), 0644)

	// Create page
	page := `<div class="hero">
    <h1>Welcome to Next.go!</h1>
    <p>A Next.js-inspired framework built in Go</p>
    <div style="margin-top: 2rem;">
        <a href="/about" style="background: #000; color: white; padding: 0.75rem 1.5rem; text-decoration: none; border-radius: 5px;">Learn More</a>
    </div>
</div>
<style>
.hero { text-align: center; padding: 4rem 0; }
.hero h1 { font-size: 3rem; margin-bottom: 1rem; }
.hero p { font-size: 1.25rem; color: #666; }
</style>`

	os.WriteFile(filepath.Join(dir, "app", "page.go.html"), []byte(page), 0644)

	// Create about page
	about := `<div class="about-page">
    <h1>About Next.go</h1>
    <p>Next.go brings Next.js concepts to Go.</p>
    <ul>
        <li>🚀 File-system routing</li>
        <li>⚡ Server-Side Rendering</li>
        <li>🔌 API Routes</li>
    </ul>
</div>`

	os.WriteFile(filepath.Join(dir, "app", "about", "page.go.html"), []byte(about), 0644)

	// Create API handler
	apiHandler := `package main

import "github.com/gin-gonic/gin"

func Handler(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "Hello from Next.go API!",
        "method":  c.Request.Method,
    })
}`

	os.WriteFile(filepath.Join(dir, "app", "api", "hello", "handler.go"), []byte(apiHandler), 0644)

	// Create package.json
	packageJSON := `{
  "name": "my-nextgo-app",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "nextgo dev",
    "build": "nextgo build",
    "start": "nextgo start"
  }
}`
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(packageJSON), 0644)

	// Create nextgo.yaml
	config := `server:
  port: 3000
  host: localhost

build:
  output: standalone
  distDir: .next
`
	os.WriteFile(filepath.Join(dir, "nextgo.yaml"), []byte(config), 0644)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
