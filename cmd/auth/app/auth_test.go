package app

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

	// Check if flags are registered
	assert.NotNil(t, cmd.Flags().Lookup("port"))
	assert.NotNil(t, cmd.Flags().Lookup("db-host"))
	assert.NotNil(t, cmd.Flags().Lookup("db-port"))
	assert.NotNil(t, cmd.Flags().Lookup("db-name"))
	assert.NotNil(t, cmd.Flags().Lookup("db-user"))
	assert.NotNil(t, cmd.Flags().Lookup("db-pass"))
}
