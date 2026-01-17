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

func TestAuthOptions_Complete(t *testing.T) {
	opts := NewAuthOptions()
	assert.NotPanics(t, func() {
		opts.Complete([]string{})
		opts.Complete([]string{"arg1", "arg2"})
	})
}

func TestAuthOptions_Config(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AuthOptions
		wantErr bool
	}{
		{
			name: "Valid config with valid database options",
			opts: &AuthOptions{
				Port: 8080,
				Database: DatabaseOptions{
					DBHost: "localhost",
					DBPort: 5432,
					DBName: "testdb",
					DBUser: "testuser",
				},
			},
			wantErr: true,
		},
		{
			name: "Config with invalid database port (triggers DSN error)",
			opts: &AuthOptions{
				Port: 8080,
				Database: DatabaseOptions{
					DBHost: "localhost",
					DBPort: 0,
					DBName: "testdb",
					DBUser: "testuser",
				},
			},
			wantErr: true,
		},
		{
			name: "Config with missing database options",
			opts: &AuthOptions{
				Port:     8080,
				Database: DatabaseOptions{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := tt.opts.Config()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				if cfg != nil {
					assert.Equal(t, tt.opts.Port, cfg.Port)
				}
			}
		})
	}
}

func TestAuthOptions_AddFlags(t *testing.T) {
	opts := NewAuthOptions()
	cmd := &cobra.Command{}
	opts.AddFlags(cmd)

	portFlag := cmd.Flags().Lookup("port")
	assert.NotNil(t, portFlag)
	assert.Equal(t, "p", portFlag.Shorthand)

	assert.NotNil(t, cmd.Flags().Lookup("db-name"))
	assert.NotNil(t, cmd.Flags().Lookup("db-host"))
}

func TestAuthOptions_Run(t *testing.T) {
	tests := []struct {
		name    string
		opts    *AuthOptions
		wantErr bool
	}{
		{
			name: "Run with invalid database port (triggers Config error)",
			opts: &AuthOptions{
				Port: 8080,
				Database: DatabaseOptions{
					DBHost: "localhost",
					DBPort: 0,
					DBName: "testdb",
					DBUser: "testuser",
				},
			},
			wantErr: true,
		},
		{
			name: "Run with missing database config",
			opts: &AuthOptions{
				Port:     8080,
				Database: DatabaseOptions{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Run()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
