package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig_ReturnsEmptyConfig(t *testing.T) {
	t.Parallel()

	cfg := NewConfig()
	require.NotNil(t, cfg)
	require.Nil(t, cfg.DB)
	require.Zero(t, cfg.Port)
}
