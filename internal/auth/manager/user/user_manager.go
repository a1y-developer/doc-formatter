package user

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/a1y/doc-formatter/internal/auth/domain/entity"
	"github.com/a1y/doc-formatter/pkg/credentials"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/copier"
)

const jwtPrivateKeyEnv = "AUTH_JWT_PRIVATE_KEY"

func loadRSAPrivateKeyFromEnv() (*rsa.PrivateKey, error) {
	pemString := os.Getenv(jwtPrivateKeyEnv)
	if pemString == "" {
		return nil, fmt.Errorf("%s environment variable is not set", jwtPrivateKeyEnv)
	}

	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM for RSA private key")
	}

	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected PEM type %q, want %q", block.Type, "PRIVATE KEY")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse PKCS#8 private key: %w", err)
	}

	key, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not an RSA key")
	}

	return key, nil
}

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

	// TODO: create new method for generate token. Now just for demo
	exp := time.Now().Add(15 * time.Minute).Unix()
	privateKey, err := loadRSAPrivateKeyFromEnv()
	if err != nil {
		return nil, 0, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   exp,
	})
	str, err := token.SignedString(privateKey)
	return &str, exp, err
}
