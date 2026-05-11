package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/nextgo/nextgo/packages/build"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the application for production",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()

		fmt.Println("Building Next.go application...")

		builder := build.New(dir)
		if err := builder.Build(); err != nil {
			fmt.Fprintf(os.Stderr, "Build error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Build complete! Output: .next/")

		// Create .next directory structure
		nextDir := filepath.Join(dir, ".next")
		dirs := []string{
			"server",
			"static",
			"server/app",
			"static/chunks",
			"cache",
		}

		for _, d := range dirs {
			os.MkdirAll(filepath.Join(nextDir, d), 0755)
		}
	},
}
