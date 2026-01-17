package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig_HasExpectedDefaults(t *testing.T) {
	t.Parallel()

	cfg := NewConfig()
	require.NotNil(t, cfg)

	require.Nil(t, cfg.DB)
	require.Equal(t, 8082, cfg.Port)
	require.Equal(t, "", cfg.EndPoint)
	require.Equal(t, "us-east-1", cfg.Region)
	require.Equal(t, "", cfg.AccessKeyID)
	require.Equal(t, "", cfg.AccessKeySecret)
	require.Equal(t, "", cfg.Bucket)
	require.False(t, cfg.ForcePathStyle)
}
