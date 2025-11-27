package constant

import "errors"

var (
	ErrEmptyEmail         = errors.New("email cannot be empty")
	ErrEmptyPassword      = errors.New("password cannot be empty")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
