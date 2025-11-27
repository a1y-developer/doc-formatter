package auth

import (
	"github.com/a1y/doc-formatter/internal/gateway/clients/auth"
)

type AuthManager struct {
	authClient auth.AuthClient
}

func NewAuthManager(authClient auth.AuthClient) *AuthManager {
	return &AuthManager{
		authClient: authClient,
	}
}
