package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCmdStorage(t *testing.T) {
	cmd := NewCmdStorage()

	assert.NotNil(t, cmd)
	assert.Equal(t, "storage", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
	assert.NotNil(t, cmd.Flags().Lookup("port"))
	assert.NotNil(t, cmd.Flags().Lookup("db-host"))
	assert.NotNil(t, cmd.Flags().Lookup("db-port"))
	assert.NotNil(t, cmd.Flags().Lookup("db-name"))
	assert.NotNil(t, cmd.Flags().Lookup("db-user"))
	assert.NotNil(t, cmd.Flags().Lookup("db-pass"))

	assert.NotNil(t, cmd.Flags().Lookup("s3-endpoint"))
	assert.NotNil(t, cmd.Flags().Lookup("s3-region"))
	assert.NotNil(t, cmd.Flags().Lookup("s3-access-key-id"))
	assert.NotNil(t, cmd.Flags().Lookup("s3-access-key-secret"))
	assert.NotNil(t, cmd.Flags().Lookup("s3-bucket"))
	assert.NotNil(t, cmd.Flags().Lookup("s3-force-path-style"))
}

func TestNewCmdStorage_RunE_Validation(t *testing.T) {
	cmd := NewCmdStorage()

	err := cmd.RunE(cmd, []string{})
	assert.Error(t, err)
}

func TestNewCmdStorage_RunE_WithFlags(t *testing.T) {
	cmd := NewCmdStorage()

	_ = cmd.Flags().Set("port", "8082")
	_ = cmd.Flags().Set("db-host", "localhost")
	_ = cmd.Flags().Set("db-port", "5432")
	_ = cmd.Flags().Set("db-name", "test_db")
	_ = cmd.Flags().Set("db-user", "test_user")
	_ = cmd.Flags().Set("db-pass", "test_pass")
	_ = cmd.Flags().Set("s3-endpoint", "http://localhost:9000")
	_ = cmd.Flags().Set("s3-bucket", "test-bucket")

	assert.NotNil(t, cmd.RunE)
}
