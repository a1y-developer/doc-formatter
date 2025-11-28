package credentials

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgon2idHash_HashPassword(t *testing.T) {
	hasher := NewDefaultArgon2idHash()
	password := "securepassword"

	t.Run("Success", func(t *testing.T) {
		hash, err := hasher.HashPassword(password, nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.True(t, strings.HasPrefix(hash, "$argon2id$"))
	})

	t.Run("SuccessWithSalt", func(t *testing.T) {
		salt := []byte("randomsalt123456")
		hash, err := hasher.HashPassword(password, salt)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.True(t, strings.HasPrefix(hash, "$argon2id$"))
	})
}

func TestCompare(t *testing.T) {
	hasher := NewDefaultArgon2idHash()
	password := "securepassword"
	wrongPassword := "wrongpassword"

	hash, err := hasher.HashPassword(password, nil)
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		match, err := Compare(password, hash)
		assert.NoError(t, err)
		assert.True(t, match)
	})

	t.Run("WrongPassword", func(t *testing.T) {
		match, err := Compare(wrongPassword, hash)
		assert.NoError(t, err)
		assert.False(t, match)
	})

	t.Run("InvalidHashFormat", func(t *testing.T) {
		invalidHash := "invalid$hash$format"
		match, err := Compare(password, invalidHash)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidHashFormat, err)
		assert.False(t, match)
	})

	t.Run("IncompatibleVersion", func(t *testing.T) {
		// Construct a hash with an incompatible version
		// $argon2id$v=18$m=65536,t=1,p=4$salt$hash
		incompatibleHash := "$argon2id$v=18$m=65536,t=1,p=4$c2FsdA$aGFzaA"
		match, err := Compare(password, incompatibleHash)
		assert.Error(t, err)
		assert.Equal(t, ErrArgon2VersionIncompatible, err)
		assert.False(t, match)
	})
}
