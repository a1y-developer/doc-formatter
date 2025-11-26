package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "auth-service",
    Short: "Auth service for AI Doc Formatter",
    Long:  "Authentication microservice with gRPC server, migrations and worker tasks.",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        panic(err)
    }
}
