package grpc

import (
    "context"
    "github.com/a1y/ai-doc-formatter/internal/auth/app"
    authpb "github.com/a1y/ai-doc-formatter/internal/shared/proto/auth"
)

type AuthServer struct {
    authpb.UnimplementedAuthServiceServer
    svc *app.AuthService
}

func NewAuthServer(s *app.AuthService) *AuthServer {
    return &AuthServer{svc: s}
}

func (s *AuthServer) Signup(ctx context.Context, req *authpb.SignupRequest) (*authpb.SignupResponse, error) {
    id, err := s.svc.Signup(ctx, req.Email, req.Password)
    if err != nil {
        return nil, err
    }
    return &authpb.SignupResponse{UserId: id}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
    token, exp, err := s.svc.Login(ctx, req.Email, req.Password)
    if err != nil {
        return nil, err
    }
    return &authpb.LoginResponse{
        AccessToken: token,
        ExpiryUnix: exp,
    }, nil
}
