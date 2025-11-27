package request

import (
	"github.com/a1y/doc-formatter/internal/gateway/domain/constant"
)

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (r *SignupRequest) Validate() error {
	if r.Email == "" {
		return constant.ErrEmptyEmail
	}
	if r.Password == "" {
		return constant.ErrEmptyPassword
	}
	return nil
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return constant.ErrEmptyEmail
	}
	if r.Password == "" {
		return constant.ErrEmptyPassword
	}
	return nil
}
