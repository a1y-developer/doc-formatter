package auth

import (
	"context"
	"log"

	authpb "github.com/a1y/doc-formatter/api/grpc/auth/v1"
	"github.com/a1y/doc-formatter/internal/gateway/domain/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient interface {
	Signup(ctx context.Context, email, password string) (*response.SignUpResponse, error)
	Login(ctx context.Context, email, password string) (*response.LoginResponse, error)
}

var _ AuthClient = &authClient{}

type authClient struct {
	conn   *grpc.ClientConn
	client authpb.AuthServiceClient
}

func NewAuthClient(addr string) AuthClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot connect to AuthService: %v", err)
	}
	c := authpb.NewAuthServiceClient(conn)
	return &authClient{
		conn:   conn,
		client: c,
	}
}
