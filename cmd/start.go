package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nextgo/nextgo/packages/server"
)
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the production server",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := os.Getwd()
		port, _ := cmd.Flags().GetString("port")

		fmt.Println("Starting Next.go production server...")

		s := server.New(dir)
		s.DevMode = false

		if port == "" {
			port = ":3000"
		} else if !strings.HasPrefix(port, ":") {
			port = ":" + port
		}

		if err := s.Start(port); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	startCmd.Flags().StringP("port", "p", "", "Port to run the production server")
}
