package cmd

import (
    "log"

    "github.com/spf13/cobra"
    gateway "github.com/a1y/ai-doc-formatter/internal/gateway"
)

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start API Gateway HTTP Server",
    Run: func(cmd *cobra.Command, args []string) {
        log.Println("Starting API Gateway on :8080")
        gateway.StartHTTPServer()
    },
}

func init() {
    rootCmd.AddCommand(serveCmd)
}
