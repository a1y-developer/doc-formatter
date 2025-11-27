package auth

import (
	"github.com/a1y/doc-formatter/internal/gateway/manager/auth"
)

type AuthHandler struct {
	authManager *auth.AuthManager
}

func NewAuthHandler(authManager *auth.AuthManager) (*AuthHandler, error) {
	return &AuthHandler{
		authManager: authManager,
	}, nil
}
