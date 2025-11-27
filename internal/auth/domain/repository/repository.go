package repository

import (
	"context"

	"github.com/a1y/doc-formatter/internal/auth/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, u *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
