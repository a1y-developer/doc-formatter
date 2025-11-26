package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/a1y/ai-doc-formatter/internal/auth/infra"
)

var migrateCmd = &cobra.Command{
    Use:   "migrate",
    Short: "Run database migrations",
    Run: func(cmd *cobra.Command, args []string) {
        db, _ := infra.NewPostgres(infra.LoadConfig())
        db.AutoMigrate(&infra.UserModel{})
        fmt.Println("Migration completed.")
    },
}

func init() {
    rootCmd.AddCommand(migrateCmd)
}
