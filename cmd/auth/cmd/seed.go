package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
    Use:   "seed",
    Short: "Seed initial data",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Running seed tasks...")
    },
}

func init() {
    rootCmd.AddCommand(seedCmd)
}
