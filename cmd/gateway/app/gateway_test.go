package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCmdGateway(t *testing.T) {
	cmd := NewCmdGateway()

	assert.NotNil(t, cmd)
	assert.Equal(t, "gateway", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)

	// Check if flags are registered
	assert.NotNil(t, cmd.Flags().Lookup("bind-address"))
	assert.NotNil(t, cmd.Flags().Lookup("auth-service"))
}
