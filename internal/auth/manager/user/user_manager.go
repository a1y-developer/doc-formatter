package user

import (
	"context"
	"time"

	"github.com/a1y/doc-formatter/internal/auth/domain/entity"
	"github.com/a1y/doc-formatter/pkg/credentials"
	"github.com/golang-jwt/jwt/v5"
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
	if !ok || err != nil {
		return nil, 0, err
	}

	// TODO: create new method for generate token. Now just for demo
	exp := time.Now().Add(15 * time.Minute).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   exp,
	})
	str, err := token.SignedString([]byte("super-secret-key"))
	return &str, exp, err
}
