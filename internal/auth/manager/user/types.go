package user

import (
	"github.com/a1y/doc-formatter/internal/auth/domain/repository"
	jwtutil "github.com/a1y/doc-formatter/internal/auth/util/jwt"
)

type UserManager struct {
	userRepo  repository.UserRepository
	jwtClaims jwtutil.TokenClaim
}

func NewUserManager(userRepo repository.UserRepository, jwtClaims jwtutil.TokenClaim) *UserManager {
	return &UserManager{
		userRepo:  userRepo,
		jwtClaims: jwtClaims,
	}
}
