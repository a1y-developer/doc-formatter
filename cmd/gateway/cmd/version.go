package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var (
    Version = "1.0.0"
    Commit  = "dev"
    BuiltAt = "local"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Show version info",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("API Gateway\nVersion: %s\nCommit: %s\nBuiltAt: %s\n", Version, Commit, BuiltAt)
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}
