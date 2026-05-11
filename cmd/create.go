package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "Create a new Next.go project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		fmt.Printf("Creating Next.go project: %s\n", projectName)

		// Create project directory
		projectDir := filepath.Join(".", projectName)
		os.MkdirAll(projectDir, 0755)

		// Initialize the project
		createProject(projectDir)

		fmt.Printf("\n✅ Success! Created %s at ./%s\n", projectName, projectName)
		fmt.Println("\nInside that directory, you can run several commands:")
		fmt.Println("  nextgo dev    - Starts the development server")
		fmt.Println("  nextgo build  - Builds the app for production")
		fmt.Println("  nextgo start  - Starts the production server")
	},
}

func createProject(dir string) {
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

	// Create layout
	layout := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.title}}</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
    </style>
</head>
<body>
    <header><nav><a href="/">Home</a> | <a href="/about">About</a></nav></header>
    <main>{{.children}}</main>
    <footer><p>Built with Next.go</p></footer>
</body>
</html>`

	os.WriteFile(filepath.Join(dir, "app", "layout.go.html"), []byte(layout), 0644)

	// Create page
	page := `<div><h1>Welcome to Next.go!</h1><p>Get started by editing app/page.go.html</p></div>`
	os.WriteFile(filepath.Join(dir, "app", "page.go.html"), []byte(page), 0644)

	// Create about page
	about := `<div><h1>About Next.go</h1><p>Next.go brings Next.js concepts to Go.</p></div>`
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

	fmt.Println("✓ Created Next.go project structure!")
}
