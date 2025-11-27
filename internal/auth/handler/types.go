package handler

import (
	authpb "github.com/a1y/doc-formatter/api/grpc/auth/v1"
	"github.com/a1y/doc-formatter/internal/auth/manager/user"
)

func NewHandler(userManager *user.UserManager) (*Handler, error) {
	return &Handler{userManager: userManager}, nil
}

type Handler struct {
	authpb.UnimplementedAuthServiceServer
	userManager *user.UserManager
}
