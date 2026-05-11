package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/nextgo/nextgo/packages/server"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the production server",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()

		fmt.Println("Starting Next.go production server...")

		s := server.New(dir)
		s.DevMode = false
		if err := s.Start(":3000"); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	},
}
