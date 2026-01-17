package request

import (
	"github.com/a1y/doc-formatter/pkg/gateway/domain/constant"
	"github.com/gin-gonic/gin"
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

func (r *SignupRequest) Decode(c *gin.Context) error {
	return decode(c, r)
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

func (r *LoginRequest) Decode(c *gin.Context) error {
	return decode(c, r)
}
