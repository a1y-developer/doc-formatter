package util

import (
	"bytes"
	"context"
	"testing"

	middleware2 "github.com/a1y/doc-formatter/pkg/gateway/middleware"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetLogger_WithAndWithoutContext(t *testing.T) {
	t.Parallel()

	logger := GetLogger(context.TODO())
	require.NotNil(t, logger)

	logger = GetLogger(context.Background())
	require.NotNil(t, logger)

	baseLogger, _ := zap.NewProduction()
	ctx := context.WithValue(context.Background(), middleware2.APILoggerKey, baseLogger)
	logger = GetLogger(ctx)
	require.Equal(t, baseLogger, logger)
}

func TestGetRunLogger_WithAndWithoutContext(t *testing.T) {
	t.Parallel()

	logger := GetRunLogger(context.TODO())
	require.NotNil(t, logger)

	logger = GetRunLogger(context.Background())
	require.NotNil(t, logger)

	baseLogger, _ := zap.NewProduction()
	ctx := context.WithValue(context.Background(), middleware2.RunLoggerKey, baseLogger)
	logger = GetRunLogger(ctx)
	require.Equal(t, baseLogger, logger)
}

func TestGetRunLoggerBuffer_WithAndWithoutContext(t *testing.T) {
	t.Parallel()

	buf := GetRunLoggerBuffer(context.TODO())
	require.NotNil(t, buf)

	buf = GetRunLoggerBuffer(context.Background())
	require.NotNil(t, buf)

	expected := &bytes.Buffer{}
	ctx := context.WithValue(context.Background(), middleware2.RunLoggerBufferKey, expected)
	buf = GetRunLoggerBuffer(ctx)
	require.Equal(t, expected, buf)
}

func TestWithRequestFields_AddsExpectedFields(t *testing.T) {
	t.Parallel()

	logger, buf := middleware2.InitLoggerBuffer("test")

	ctx := context.WithValue(context.Background(), middleware2.APILoggerKey, logger)

	enriched := WithRequestFields(ctx, "req-1", "user-1", "trace-1")
	enriched.Info("message")
	_ = enriched.Sync()

	out := buf.String()
	require.Contains(t, out, "req-1")
	require.Contains(t, out, "user-1")
	require.Contains(t, out, "trace-1")
}
