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
	// Complete is a no-op, but we should test it doesn't panic
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
			wantErr: true, // Will fail because DB is not available
		},
		{
			name: "Config with invalid database port (triggers DSN error)",
			opts: &AuthOptions{
				Port: 8080,
				Database: DatabaseOptions{
					DBHost: "localhost",
					DBPort: 0, // Invalid port - will cause DSN parsing error
					DBName: "testdb",
					DBUser: "testuser",
				},
			},
			wantErr: true, // Will fail due to invalid DSN (port=0)
		},
		{
			name: "Config with missing database options",
			opts: &AuthOptions{
				Port:     8080,
				Database: DatabaseOptions{
					// Missing required fields
				},
			},
			wantErr: true, // Will fail due to invalid DSN
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := tt.opts.Config()
			if tt.wantErr {
				assert.Error(t, err)
				// When Config returns an error, cfg should be nil (line 39)
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

	// Check if port flag is added
	portFlag := cmd.Flags().Lookup("port")
	assert.NotNil(t, portFlag)
	assert.Equal(t, "p", portFlag.Shorthand)

	// Check if database flags are added (indirectly verifying DatabaseOptions.AddFlags is called)
	assert.NotNil(t, cmd.Flags().Lookup("db-name"))
	assert.NotNil(t, cmd.Flags().Lookup("db-host"))

	// Test that AddFlags handles invalid PortEnv gracefully
	// Note: PortEnv is a package-level variable set at init time,
	// so we can't change it in tests, but the code path for invalid PortEnv
	// is covered when PortEnv is invalid (defaults to DefaultPort)
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
					DBPort: 0, // Invalid port - will cause Config() to fail
					DBName: "testdb",
					DBUser: "testuser",
				},
			},
			wantErr: true, // Will fail at Config() due to invalid database options (line 57-58)
		},
		{
			name: "Run with missing database config",
			opts: &AuthOptions{
				Port:     8080,
				Database: DatabaseOptions{
					// Missing required fields - will cause Config to fail
				},
			},
			wantErr: true, // Will fail at Config() due to invalid database options
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run() will fail early at Config() when database options are invalid
			// This tests the error path in Run() at lines 57-58
			err := tt.opts.Run()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				// Note: Run() starts a gRPC server which runs indefinitely
				// So we can't test the success path without mocking or using a timeout
				assert.NoError(t, err)
			}
		})
	}
}
