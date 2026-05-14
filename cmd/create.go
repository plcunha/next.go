package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Default project templates embedded at build time.
//
//go:embed templates/layout.go.html
var defaultLayout string

//go:embed templates/home.go.html
var defaultHome string

//go:embed templates/about.go.html
var defaultAbout string

//go:embed templates/api_handler.go
var defaultAPIHandler string

//go:embed templates/nextgo.yaml
var defaultConfig string

//go:embed templates/package.json
var defaultPackageJSON string

var createCmd = &cobra.Command{
	Use:   "create [project-name]",
	Short: "Create a new Next.go project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		fmt.Printf("Creating Next.go project: %s\n", projectName)

		projectDir := filepath.Join(".", projectName)
		os.MkdirAll(projectDir, 0755)

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

	// Write embedded templates
	os.WriteFile(filepath.Join(dir, "app", "layout.go.html"), []byte(defaultLayout), 0644)
	os.WriteFile(filepath.Join(dir, "app", "page.go.html"), []byte(defaultHome), 0644)
	os.WriteFile(filepath.Join(dir, "app", "about", "page.go.html"), []byte(defaultAbout), 0644)
	os.WriteFile(filepath.Join(dir, "app", "api", "hello", "handler.go"), []byte(defaultAPIHandler), 0644)
	os.WriteFile(filepath.Join(dir, "nextgo.yaml"), []byte(defaultConfig), 0644)
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(defaultPackageJSON), 0644)

	fmt.Println("✓ Created Next.go project structure!")
}
