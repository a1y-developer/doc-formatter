package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCmdAuth(t *testing.T) {
	cmd := NewCmdAuth()

	assert.NotNil(t, cmd)
	assert.Equal(t, "auth", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
	assert.NotNil(t, cmd.Flags().Lookup("port"))
	assert.NotNil(t, cmd.Flags().Lookup("db-host"))
	assert.NotNil(t, cmd.Flags().Lookup("db-port"))
	assert.NotNil(t, cmd.Flags().Lookup("db-name"))
	assert.NotNil(t, cmd.Flags().Lookup("db-user"))
	assert.NotNil(t, cmd.Flags().Lookup("db-pass"))
}

func TestNewCmdAuth_RunE_Validation(t *testing.T) {
	cmd := NewCmdAuth()

	err := cmd.RunE(cmd, []string{})
	assert.Error(t, err)
}

func TestNewCmdAuth_RunE_WithFlags(t *testing.T) {
	cmd := NewCmdAuth()

	cmd.Flags().Set("port", "8080")
	cmd.Flags().Set("db-host", "localhost")
	cmd.Flags().Set("db-port", "5432")
	cmd.Flags().Set("db-name", "test_db")
	cmd.Flags().Set("db-user", "test_user")
	cmd.Flags().Set("db-pass", "test_pass")
	assert.NotNil(t, cmd.RunE)
}
