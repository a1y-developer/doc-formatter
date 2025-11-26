package infra

import "golang.org/x/crypto/bcrypt"

type PasswordHasher struct{}

func NewPasswordHasher() *PasswordHasher { return &PasswordHasher{} }

func (p *PasswordHasher) Hash(pw string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
    return string(bytes), err
}

func (p *PasswordHasher) Compare(hashed, plain string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
