package middleware

import (
	"bytes"
	"context"
	"os"
	"time"

	"github.com/a1y/doc-formatter/internal/gateway"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerOptions is a structured configuration for building a zap logger
// in a way that matches common big-tech practices.
type LoggerOptions struct {
	// ServiceName is the logical name of the service emitting logs.
	ServiceName string
	// Environment is typically "dev" or "prod".
	Environment string
	// Level is the minimum log level: debug, info, warn, error.
	Level string
	// Format is either "json" (recommended for production) or "console".
	Format string
	// FilePath is the optional log file path. If empty, logs go only to stdout/stderr.
	FilePath string
	// MaxSize is the maximum size in megabytes of the log file before it gets rotated.
	MaxSize int
	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int
	// MaxAge is the maximum number of days to retain old log files.
	MaxAge int
	// Compress determines if the rotated log files should be compressed.
	Compress bool
	// Sample enables log sampling in production to reduce volume.
	Sample bool
}

// InitLoggerWithOptions builds a zap.Logger using LoggerOptions, logging to
// both stdout and an optional rotating file, with JSON or console encoding
// and optional sampling. This follows patterns used in large-scale Go systems.
func InitLoggerWithOptions(opts LoggerOptions) (*zap.Logger, error) {
	env := opts.Environment
	if env == "" {
		env = "dev"
	}

	level := zapcore.InfoLevel
	if opts.Level != "" {
		if err := level.Set(opts.Level); err != nil {
			// Fall back to info if the level string is invalid.
			level = zapcore.InfoLevel
		}
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
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

	newEncoder := func() zapcore.Encoder {
		if opts.Format == "console" && env == "dev" {
			return zapcore.NewConsoleEncoder(encoderCfg)
		}
		return zapcore.NewJSONEncoder(encoderCfg)
	}

	// Always log to stdout so container/orchestration layers can collect logs.
	stdoutSyncer := zapcore.AddSync(os.Stdout)
	cores := []zapcore.Core{
		zapcore.NewCore(newEncoder(), stdoutSyncer, level),
	}

	// Optionally log to a rotating file for on-prem deployments.
	if opts.FilePath != "" {
		maxSize := opts.MaxSize
		if maxSize <= 0 {
			maxSize = 100
		}
		maxBackups := opts.MaxBackups
		if maxBackups < 0 {
			maxBackups = 0
		}
		maxAge := opts.MaxAge
		if maxAge < 0 {
			maxAge = 0
		}

		fileLogger := &lumberjack.Logger{
			Filename:   opts.FilePath,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   opts.Compress,
		}

		fileSyncer := zapcore.AddSync(fileLogger)
		cores = append(cores, zapcore.NewCore(newEncoder(), fileSyncer, level))
	}

	core := zapcore.NewTee(cores...)

	// Enable sampling in production to reduce log volume while retaining detail.
	if opts.Sample && env == "prod" {
		core = zapcore.NewSamplerWithOptions(core, time.Second, 100, 10)
	}

	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	fields := []zap.Field{
		zap.String("env", env),
	}
	if opts.ServiceName != "" {
		fields = append(fields, zap.String("service", opts.ServiceName))
	}

	return logger.With(fields...), nil
}

// InitLoggerBuffer initializes a zap logger that writes JSON logs
// into an in-memory buffer. This is useful for tests or temporary
// logging sinks.
func InitLoggerBuffer(name string) (*zap.Logger, *bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(buf),
		zap.InfoLevel,
	)

	logger := zap.New(core)
	if name != "" {
		logger = logger.Named(name)
	}

	return logger, buf, nil
}

// APILoggerMiddleware initializes a zap logger using the gateway logging config,
// attaches it to the request context, and emits structured JSON logs
// for each HTTP request with big-tech-style fields.
func APILoggerMiddleware(cfg gateway.LoggingConfig, serviceName string) gin.HandlerFunc {
	logger, err := InitLoggerWithOptions(LoggerOptions{
		ServiceName: serviceName,
		Environment: cfg.Environment,
		Level:       cfg.Level,
		Format:      cfg.Format,
		FilePath:    cfg.FilePath,
		MaxSize:     cfg.MaxSize,
		MaxBackups:  cfg.MaxBackups,
		MaxAge:      cfg.MaxAge,
		Compress:    cfg.Compress,
		Sample:      cfg.Sample,
	})
	if err != nil {
		// Fail fast on startup misconfiguration.
		panic(err)
	}

	return func(c *gin.Context) {
		start := time.Now()

		req := c.Request

		// Request ID: prefer incoming header, otherwise generate one and echo back.
		requestID := req.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Writer.Header().Set("X-Request-Id", requestID)
		}

		traceID := req.Header.Get("X-Trace-Id")
		spanID := req.Header.Get("X-Span-Id")

		reqLogger := logger.With(
			zap.String("request_id", requestID),
			zap.String("trace", traceID),
			zap.String("span", spanID),
			zap.String("http_method", req.Method),
			zap.String("http_path", req.URL.Path),
			zap.String("remote_ip", c.ClientIP()),
			zap.String("user_agent", req.UserAgent()),
		)

		// Attach the request-scoped logger to the context so downstream code can use it.
		reqCtx := context.WithValue(req.Context(), APILoggerKey, reqLogger)
		c.Request = c.Request.WithContext(reqCtx)

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
