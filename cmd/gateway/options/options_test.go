package options

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewOptions(t *testing.T) {
	opts := NewOptions()
	assert.NotNil(t, opts)
}

func TestOptions_AddFlags(t *testing.T) {
	opts := NewOptions()
	cmd := &cobra.Command{}
	opts.AddFlags(cmd)

	// Check if flags are added
	assert.NotNil(t, cmd.Flags().Lookup("bind-address"))
	assert.NotNil(t, cmd.Flags().Lookup("auth-service"))
}

func TestOptions_Config(t *testing.T) {
	opts := &Options{
		Address:     ":9090",
		AuthService: ":9091",
	}
	cfg, err := opts.Config()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, ":9090", cfg.Address)
	assert.Equal(t, ":9091", cfg.AuthService)
}

func TestOptions_Validate(t *testing.T) {
	opts := NewOptions()
	err := opts.Validate()
	assert.NoError(t, err)
}
