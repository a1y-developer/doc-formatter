package middleware

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/a1y/doc-formatter/pkg/gateway"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/klog/v2"
)

func InitLogger(logFilePath string, name string) *zap.Logger {
	var syncer zapcore.WriteSyncer

	if logFilePath == "" {
		syncer = zapcore.AddSync(os.Stdout)
	} else {
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			// if directory does not exist, try to create the directory
			if os.IsNotExist(err) {
				logFileParent := filepath.Dir(logFilePath)
				klog.Infof("Log directory does not exist, trying to create the directory in %s", logFileParent)
				err = os.MkdirAll(logFileParent, 0o755)
				if err != nil {
					klog.Fatalf("Failed to create log directory: %v", err)
				}
				// Try opening the file again after creating directory
				logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
				if err != nil {
					klog.Fatalf("Failed to open log file: %v", err)
				}
			} else {
				klog.Fatalf("Failed to open log file: %v", err)
			}
		}
		syncer = zapcore.AddSync(logFile)
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		syncer,
		zap.InfoLevel,
	)

	return zap.New(core).Named(name)
}

func InitLoggerBuffer(name string) (*zap.Logger, *bytes.Buffer) {
	var buffer bytes.Buffer

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(&buffer),
		zap.InfoLevel,
	)

	return zap.New(core).Named(name), &buffer
}

// APILoggerMiddleware injects a logger, configured with a request ID,
// into the request context for use throughout the request's lifecycle.
func APILoggerMiddleware(cfg gateway.LoggingConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		req := c.Request
		ctx := req.Context()

		requestID := req.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Writer.Header().Set("X-Request-Id", requestID)
		}

		traceID := GetTraceID(ctx)

		loggerName := traceID
		if loggerName == "" {
			loggerName = "gateway"
		}

		// Use configured file path
		logger := InitLogger(cfg.FilePath, loggerName)
		runLogger, logBuffer := InitLoggerBuffer(loggerName)

		reqLogger := logger.With(
			zap.String("request_id", requestID),
			zap.String("trace", traceID),
			zap.String("http_method", req.Method),
			zap.String("http_path", req.URL.Path),
			zap.String("remote_ip", c.ClientIP()),
			zap.String("user_agent", req.UserAgent()),
		)

		ctx = context.WithValue(ctx, APILoggerKey, reqLogger)
		ctx = context.WithValue(ctx, RunLoggerKey, runLogger)
		ctx = context.WithValue(ctx, RunLoggerBufferKey, logBuffer)

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()

		reqLogger.Info("request completed",
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.Int("response_size", size),
		)
	}
}
