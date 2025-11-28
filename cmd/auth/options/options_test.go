package options

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthOptions(t *testing.T) {
	opts := NewAuthOptions()
	assert.NotNil(t, opts)
	assert.Equal(t, DefaultPort, opts.Port)
	assert.NotNil(t, opts.Database)
}

func TestAuthOptions_Validate(t *testing.T) {
	opts := NewAuthOptions()
	err := opts.Validate()
	assert.NoError(t, err)
}

func TestAuthOptions_AddFlags(t *testing.T) {
	opts := NewAuthOptions()
	cmd := &cobra.Command{}
	opts.AddFlags(cmd)

	// Check if port flag is added
	portFlag := cmd.Flags().Lookup("port")
	assert.NotNil(t, portFlag)
	assert.Equal(t, "p", portFlag.Shorthand)

	// Check if database flags are added (indirectly verifying DatabaseOptions.AddFlags is called)
	assert.NotNil(t, cmd.Flags().Lookup("db-name"))
	assert.NotNil(t, cmd.Flags().Lookup("db-host"))
}
