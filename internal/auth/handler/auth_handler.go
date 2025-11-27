package handler

import (
	"context"

	authpb "github.com/a1y/doc-formatter/api/grpc/auth/v1"
	"github.com/a1y/doc-formatter/internal/auth/domain/entity"
)

func (h *Handler) Signup(ctx context.Context, req *authpb.SignupRequest) (*authpb.SignupResponse, error) {
	userEntity := entity.User{
		Email:    req.Email,
		Password: req.Password,
	}
	userResponse, err := h.userManager.CreateUser(ctx, &userEntity)
	if err != nil {
		return nil, err
	}
	return &authpb.SignupResponse{
		UserId: userResponse.ID.String(),
	}, nil
}

func (h *Handler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	token, exp, err := h.userManager.LoginUser(ctx, &entity.User{Email: req.Email, Password: req.Password})
	if err != nil || token == nil {
		return nil, err
	}
	return &authpb.LoginResponse{
		AccessToken: *token,
		ExpiryUnix:  exp,
	}, nil
}
