package entity

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `yaml:"id" json:"id"`
	Email      string    `yaml:"email" json:"email"`
	Password   string    `yaml:"password" json:"password"`
	IsVerified bool      `yaml:"is_verified" json:"is_verified"`
}

func (u *User) Validate() error {
	if u.ID == uuid.Nil {
		return errors.New("id is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
