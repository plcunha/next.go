package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nextgo",
	Short: "Next.go - The React Framework in Go",
	Long: `Next.go is a Go implementation of Next.js framework concepts.
It provides file-system routing, SSR, SSG, API routes, and more.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(createCmd)
}
