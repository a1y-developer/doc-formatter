package user

import (
	"context"
	"errors"
	"time"

	"github.com/a1y/doc-formatter/internal/auth/domain/entity"
	"github.com/a1y/doc-formatter/pkg/credentials"
	"github.com/jinzhu/copier"
)

func (u *UserManager) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	var createdEntity entity.User
	if err := copier.Copy(&createdEntity, &user); err != nil {
		return nil, err
	}
	argon2iHash := credentials.NewDefaultArgon2idHash()
	hashedPassword, err := argon2iHash.HashPassword(createdEntity.Password, nil)
	if err != nil {
		return nil, err
	}
	createdEntity.Password = hashedPassword

	if err := u.userRepo.Create(ctx, &createdEntity); err != nil {
		return nil, err
	}
	return &createdEntity, nil
}

func (u *UserManager) LoginUser(ctx context.Context, userEntity *entity.User) (*string, int64, error) {
	user, err := u.userRepo.GetByEmail(ctx, userEntity.Email)
	if err != nil {
		return nil, 0, err
	}
	ok, err := credentials.Compare(userEntity.Password, user.Password)
	if err != nil {
		return nil, 0, err
	}
	if !ok {
		return nil, 0, errors.New("invalid credentials")
	}

	tokenString, exp, err := u.jwtClaims.GenerateToken(user.ID, user.Email, 15*time.Minute)
	if err != nil {
		return nil, 0, err
	}

	return &tokenString, exp, nil
}
