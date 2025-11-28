package options

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		options *DatabaseOptions
		wantErr bool
	}{
		{
			name: "Valid options",
			options: &DatabaseOptions{
				DBHost: "localhost",
				DBName: "testdb",
				DBUser: "user",
				DBPort: 5432,
			},
			wantErr: false,
		},
		{
			name: "Missing DBHost",
			options: &DatabaseOptions{
				DBName: "testdb",
				DBUser: "user",
				DBPort: 5432,
			},
			wantErr: true,
		},
		{
			name: "Missing DBName",
			options: &DatabaseOptions{
				DBHost: "localhost",
				DBUser: "user",
				DBPort: 5432,
			},
			wantErr: true,
		},
		{
			name: "Missing DBUser",
			options: &DatabaseOptions{
				DBHost: "localhost",
				DBName: "testdb",
				DBPort: 5432,
			},
			wantErr: true,
		},
		{
			name: "Missing DBPort",
			options: &DatabaseOptions{
				DBHost: "localhost",
				DBName: "testdb",
				DBUser: "user",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatabaseOptions_AddFlags(t *testing.T) {
	opts := &DatabaseOptions{}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	opts.AddFlags(fs)

	// Check if flags are added
	assert.NotNil(t, fs.Lookup("db-name"))
	assert.NotNil(t, fs.Lookup("db-user"))
	assert.NotNil(t, fs.Lookup("db-pass"))
	assert.NotNil(t, fs.Lookup("db-host"))
	assert.NotNil(t, fs.Lookup("db-port"))
	assert.NotNil(t, fs.Lookup("auto-migrate"))

	// Verify default values (assuming env vars are not set or set to defaults in the environment where tests run)
	// Note: Since types.go reads from os.Getenv, the default values depend on the environment.
	// We can at least check that they have *some* value or the expected type.

	// For a more robust test, we could set/unset env vars, but that might interfere with parallel tests.
	// Here we just check availability.
}
