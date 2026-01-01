package gateway

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig_HasDefaultLoggingConfig(t *testing.T) {
	t.Parallel()

	cfg := NewConfig()
	require.NotNil(t, cfg)

	require.Equal(t, "info", cfg.Logging.Level)
	require.Equal(t, "json", cfg.Logging.Format)
	require.Equal(t, "", cfg.Logging.FilePath)
	require.Equal(t, 100, cfg.Logging.MaxSize)
	require.Equal(t, 3, cfg.Logging.MaxBackups)
	require.Equal(t, 7, cfg.Logging.MaxAge)
	require.True(t, cfg.Logging.Compress)
	require.Equal(t, "dev", cfg.Logging.Environment)
	require.False(t, cfg.Logging.Sample)
}
