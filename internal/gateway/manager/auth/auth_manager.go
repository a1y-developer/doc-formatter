package auth

import (
	"context"

	"github.com/a1y/doc-formatter/internal/gateway/domain/request"
	"github.com/a1y/doc-formatter/internal/gateway/domain/response"
)

func (m *AuthManager) Signup(ctx context.Context, request request.SignupRequest) (*response.SignUpResponse, error) {
	return m.authClient.Signup(ctx, request.Email, request.Password)
}

func (m *AuthManager) Login(ctx context.Context, request request.LoginRequest) (*response.LoginResponse, error) {
	return m.authClient.Login(ctx, request.Email, request.Password)
}
