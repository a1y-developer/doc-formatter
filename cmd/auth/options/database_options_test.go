package options

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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
	tests := []struct {
		name           string
		setDBPortEnv   bool
		dbPortEnvVal   string
		setAutoMigrate bool
		autoMigrateVal string
		expectDBPort   int
		expectAutoMig  bool
	}{
		{
			name:           "AddFlags with valid env vars",
			setDBPortEnv:   true,
			dbPortEnvVal:   "3306",
			setAutoMigrate: true,
			autoMigrateVal: "true",
			expectDBPort:   3306,
			expectAutoMig:  true,
		},
		{
			name:           "AddFlags with invalid DBPortEnv",
			setDBPortEnv:   true,
			dbPortEnvVal:   "invalid",
			setAutoMigrate: false,
			expectDBPort:   DefaultDBPort,
			expectAutoMig:  false,
		},
		{
			name:           "AddFlags with invalid AutoMigrateEnv",
			setDBPortEnv:   false,
			setAutoMigrate: true,
			autoMigrateVal: "invalid",
			expectDBPort:   DefaultDBPort,
			expectAutoMig:  false,
		},
		{
			name:           "AddFlags without env vars",
			setDBPortEnv:   false,
			setAutoMigrate: false,
			expectDBPort:   DefaultDBPort,
			expectAutoMig:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalDBPortEnv := os.Getenv("AUTH_DB_PORT")
			originalAutoMigrateEnv := os.Getenv("AUTH_AUTO_MIGRATE")
			defer func() {
				os.Setenv("AUTH_DB_PORT", originalDBPortEnv)
				os.Setenv("AUTH_AUTO_MIGRATE", originalAutoMigrateEnv)
			}()

			if tt.setDBPortEnv {
				os.Setenv("AUTH_DB_PORT", tt.dbPortEnvVal)
			} else {
				os.Unsetenv("AUTH_DB_PORT")
			}

			if tt.setAutoMigrate {
				os.Setenv("AUTH_AUTO_MIGRATE", tt.autoMigrateVal)
			} else {
				os.Unsetenv("AUTH_AUTO_MIGRATE")
			}

			opts := &DatabaseOptions{}
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			opts.AddFlags(fs)

			assert.NotNil(t, fs.Lookup("db-name"))
			assert.NotNil(t, fs.Lookup("db-user"))
			assert.NotNil(t, fs.Lookup("db-pass"))
			assert.NotNil(t, fs.Lookup("db-host"))
			assert.NotNil(t, fs.Lookup("db-port"))
			assert.NotNil(t, fs.Lookup("auto-migrate"))
		})
	}
}

func TestDatabaseOptions_InstallDB(t *testing.T) {
	tests := []struct {
		name string
		opts *DatabaseOptions
	}{
		{
			name: "InstallDB with valid options",
			opts: &DatabaseOptions{
				DBHost:     "localhost",
				DBPort:     5432,
				DBName:     "testdb",
				DBUser:     "testuser",
				DBPassword: "testpass",
			},
		},
		{
			name: "InstallDB with missing host",
			opts: &DatabaseOptions{
				DBPort: 5432,
				DBName: "testdb",
				DBUser: "testuser",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := tt.opts.InstallDB()
			_ = db
			_ = err
		})
	}
}

func TestDatabaseOptions_ApplyTo(t *testing.T) {
	tests := []struct {
		name        string
		opts        *DatabaseOptions
		autoMigrate bool
		wantErr     bool
	}{
		{
			name: "ApplyTo with valid options without auto-migrate",
			opts: &DatabaseOptions{
				DBHost:     "localhost",
				DBPort:     5432,
				DBName:     "testdb",
				DBUser:     "testuser",
				DBPassword: "testpass",
			},
			autoMigrate: false,
			wantErr:     true,
		},
		{
			name: "ApplyTo with valid options with auto-migrate",
			opts: &DatabaseOptions{
				DBHost:     "localhost",
				DBPort:     5432,
				DBName:     "testdb",
				DBUser:     "testuser",
				DBPassword: "testpass",
			},
			autoMigrate: true,
			wantErr:     true,
		},
		{
			name:        "ApplyTo with invalid options",
			opts:        &DatabaseOptions{},
			autoMigrate: false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.opts.AutoMigrate = tt.autoMigrate
			var db *gorm.DB
			err := tt.opts.ApplyTo(&db)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}
		})
	}
}
