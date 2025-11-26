package cmd

import (
    "log"
    "net"

    "github.com/spf13/cobra"

    "github.com/a1y/ai-doc-formatter/internal/auth/infra"
    "github.com/a1y/ai-doc-formatter/internal/auth/app"
    authgrpc "github.com/a1y/ai-doc-formatter/internal/auth/transport/grpc"
    authpb "github.com/a1y/ai-doc-formatter/internal/shared/proto/auth"
    "google.golang.org/grpc"
)

var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the AuthService gRPC server",
    Run: func(cmd *cobra.Command, args []string) {
        cfg := infra.LoadConfig()
        db, _ := infra.NewPostgres(cfg)

        repo := infra.NewUserRepository(db)
        pw := infra.NewPasswordHasher()
        jwt := infra.NewJWTGenerator("super-secret-key")

        svc := app.NewAuthService(repo, jwt, pw)

        lis, _ := net.Listen("tcp", ":9000")
        server := grpc.NewServer()

        authpb.RegisterAuthServiceServer(server, authgrpc.NewAuthServer(svc))

        log.Println("AuthService running at :9000")
        server.Serve(lis)
    },
}

func init() {
    rootCmd.AddCommand(serveCmd)
}
