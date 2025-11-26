package infra

import (
    "time"
    "github.com/a1y/ai-doc-formatter/internal/auth/domain"
    "github.com/golang-jwt/jwt/v5"
)

type JWTGenerator struct {
    secret []byte
}

func NewJWTGenerator(secret string) *JWTGenerator {
    return &JWTGenerator{secret: []byte(secret)}
}

func (j *JWTGenerator) GenerateAccessToken(u *domain.User) (string, int64, error) {
    exp := time.Now().Add(15 * time.Minute).Unix()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": u.ID,
        "email": u.Email,
        "exp": exp,
    })
    str, err := token.SignedString(j.secret)
    return str, exp, err
}
