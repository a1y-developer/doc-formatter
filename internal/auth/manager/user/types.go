package user

import (
	"github.com/a1y/doc-formatter/internal/auth/domain/repository"
)

type UserManager struct {
	userRepo repository.UserRepository
}

func NewUserManager(userRepo repository.UserRepository) *UserManager {
	return &UserManager{
		userRepo: userRepo,
	}
}
