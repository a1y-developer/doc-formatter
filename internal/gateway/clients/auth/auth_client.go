package auth

import (
	"context"
	"time"

	authpb "github.com/a1y/doc-formatter/api/grpc/auth/v1"
	"github.com/a1y/doc-formatter/internal/gateway/domain/response"
)

func (a *authClient) Signup(ctx context.Context, email, password string) (*response.SignUpResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := a.client.Signup(ctx, &authpb.SignupRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return &response.SignUpResponse{UserID: resp.GetUserId()}, nil
}

func (a *authClient) Login(ctx context.Context, email, password string) (*response.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := a.client.Login(ctx, &authpb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return &response.LoginResponse{
		AccessToken: resp.GetAccessToken(),
		ExpiryUnix:  resp.GetExpiryUnix(),
	}, nil
}
