package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUser_Validate_Success(t *testing.T) {
	t.Parallel()

	u := &User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "secret",
	}

	require.NoError(t, u.Validate())
}

func TestUser_Validate_MissingEmail(t *testing.T) {
	t.Parallel()

	u := &User{
		Password: "secret",
	}

	err := u.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "email is required")
}

func TestUser_Validate_MissingPassword(t *testing.T) {
	t.Parallel()

	u := &User{
		Email: "test@example.com",
	}

	err := u.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "password is required")
}
