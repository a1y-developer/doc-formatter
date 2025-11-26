package app

import (
    "context"
    "github.com/a1y/ai-doc-formatter/internal/auth/domain"
    "github.com/google/uuid"
)

type UserRepository interface {
    Create(ctx context.Context, u *domain.User) error
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type TokenGenerator interface {
    GenerateAccessToken(u *domain.User) (string, int64, error)
}

type PasswordHasher interface {
    Hash(pw string) (string, error)
    Compare(hashed, plain string) bool
}

type AuthService struct {
    users UserRepository
    jwt   TokenGenerator
    pw    PasswordHasher
}

func NewAuthService(u UserRepository, j TokenGenerator, p PasswordHasher) *AuthService {
    return &AuthService{users: u, jwt: j, pw: p}
}

func (a *AuthService) Signup(ctx context.Context, email, password string) (string, error) {
    if existing, _ := a.users.FindByEmail(ctx, email); existing != nil {
        return "", domain.ErrEmailExists
    }

    hash, _ := a.pw.Hash(password)

    u := &domain.User{
        ID:           uuid.NewString(),
        Email:        email,
        Password:     hash,
    }

    if err := a.users.Create(ctx, u); err != nil {
        return "", err
    }

    return u.ID, nil
}

func (a *AuthService) Login(ctx context.Context, email, password string) (string, int64, error) {
    user, err := a.users.FindByEmail(ctx, email)
    if err != nil || !a.pw.Compare(user.Password, password) {
        return "", 0, domain.ErrInvalidCredentials
    }

    return a.jwt.GenerateAccessToken(user)
}
