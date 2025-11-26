package clients

import (
    "context"
    "log"
    "time"

    authpb "github.com/a1y/ai-doc-formatter/internal/shared/proto/auth"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
    conn   *grpc.ClientConn
    client authpb.AuthServiceClient
}

func NewAuthClient(addr string) *AuthClient {
    // addr: "localhost:9000" (AuthService)
    conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("cannot connect to AuthService: %v", err)
    }
    c := authpb.NewAuthServiceClient(conn)
    return &AuthClient{
        conn:   conn,
        client: c,
    }
}

func (a *AuthClient) Signup(ctx context.Context, email, password string) (*authpb.SignupResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    return a.client.Signup(ctx, &authpb.SignupRequest{
        Email:    email,
        Password: password,
    })
}

func (a *AuthClient) Login(ctx context.Context, email, password string) (*authpb.LoginResponse, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    return a.client.Login(ctx, &authpb.LoginRequest{
        Email:    email,
        Password: password,
    })
}
