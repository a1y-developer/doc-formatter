package credentials

import "errors"

var (
	ErrInvalidHashFormat         = errors.New("invalid hash format")
	ErrArgon2VersionIncompatible = errors.New("argon2 version incompatible")
)
