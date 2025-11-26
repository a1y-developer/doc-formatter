package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "gateway",
    Short: "API Gateway for AI Doc Formatter",
    Long:  "API Gateway for routing HTTP requests to internal microservices via gRPC",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        panic(err)
    }
}
