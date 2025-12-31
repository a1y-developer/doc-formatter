package options

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewStorageOptions(t *testing.T) {
	opts := NewStorageOptions()
	assert.NotNil(t, opts)
	assert.Equal(t, DefaultPort, opts.Port)
	assert.NotNil(t, opts.Database)
}

func TestStorageOptions_Validate(t *testing.T) {
	opts := &StorageOptions{
		Database: DatabaseOptions{
			DBHost: "localhost",
			DBName: "testdb",
			DBUser: "user",
			DBPort: 5432,
		},
	}
	err := opts.Validate()
	assert.NoError(t, err)
}

func TestStorageOptions_Complete(t *testing.T) {
	opts := NewStorageOptions()
	assert.NotPanics(t, func() {
		opts.Complete([]string{})
		opts.Complete([]string{"arg1", "arg2"})
	})
}

func TestStorageOptions_Config(t *testing.T) {
	tests := []struct {
		name    string
		opts    *StorageOptions
		wantErr bool
	}{
		{
			name: "Valid config with valid database options",
			opts: &StorageOptions{
				Port: 8082,
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
			name: "Config with invalid database port",
			opts: &StorageOptions{
				Port: 8082,
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
			opts: &StorageOptions{
				Port:     8082,
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

func TestStorageOptions_AddFlags(t *testing.T) {
	opts := NewStorageOptions()
	cmd := &cobra.Command{}
	opts.AddFlags(cmd)

	// Check if port flag is added.
	portFlag := cmd.Flags().Lookup("port")
	assert.NotNil(t, portFlag)
	assert.Equal(t, "p", portFlag.Shorthand)

	assert.NotNil(t, cmd.Flags().Lookup("s3-endpoint"))
	assert.NotNil(t, cmd.Flags().Lookup("s3-region"))
	assert.NotNil(t, cmd.Flags().Lookup("s3-access-key-id"))
	assert.NotNil(t, cmd.Flags().Lookup("s3-access-key-secret"))
	assert.NotNil(t, cmd.Flags().Lookup("s3-bucket"))
	assert.NotNil(t, cmd.Flags().Lookup("s3-force-path-style"))

	assert.NotNil(t, cmd.Flags().Lookup("db-name"))
	assert.NotNil(t, cmd.Flags().Lookup("db-host"))
}

func TestStorageOptions_Run(t *testing.T) {
	tests := []struct {
		name    string
		opts    *StorageOptions
		wantErr bool
	}{
		{
			name: "Run with invalid database port (triggers Config error)",
			opts: &StorageOptions{
				Port: 8082,
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
			opts: &StorageOptions{
				Port:     8082,
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
