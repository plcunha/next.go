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

		// Check if this looks like a Next.go project
		if !fileExists(filepath.Join(dir, "nextgo.yaml")) && !fileExists(filepath.Join(dir, "app")) {
			fmt.Println("No Next.go project found. Run 'nextgo create <name>' first.")
			os.Exit(1)
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
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
