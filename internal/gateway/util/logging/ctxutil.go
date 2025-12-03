package util

import (
	"bytes"
	"context"

	"github.com/a1y/doc-formatter/internal/gateway/middleware"
	"go.uber.org/zap"
)

// GetLogger returns the API logger from the given context.
func GetLogger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.NewNop()
	}

	if logger, ok := ctx.Value(middleware.APILoggerKey).(*zap.Logger); ok && logger != nil {
		return logger
	}

	return zap.NewNop()
}

// GetRunLogger returns the run logger from the given context.
func GetRunLogger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.NewNop()
	}

	if logger, ok := ctx.Value(middleware.RunLoggerKey).(*zap.Logger); ok && logger != nil {
		return logger
	}

	return zap.NewNop()
}

// GetRunLoggerBuffer returns the run logger buffer from the given context.
func GetRunLoggerBuffer(ctx context.Context) *bytes.Buffer {
	if ctx == nil {
		return &bytes.Buffer{}
	}

	if buffer, ok := ctx.Value(middleware.RunLoggerBufferKey).(*bytes.Buffer); ok && buffer != nil {
		return buffer
	}

	return &bytes.Buffer{}
}

// WithRequestFields returns a logger enriched with standard request-related fields.
// This is useful in handlers or services that want to log with consistent
// metadata such as request ID, user ID, and trace ID.
func WithRequestFields(ctx context.Context, requestID, userID, traceID string) *zap.Logger {
	logger := GetLogger(ctx)

	fields := make([]zap.Field, 0, 3)
	if requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}
	if userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}
	if traceID != "" {
		fields = append(fields, zap.String("trace", traceID))
	}

	if len(fields) == 0 {
		return logger
	}

	return logger.With(fields...)
}
