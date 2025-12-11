package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/a1y/doc-formatter/internal/gateway"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TestInitLoggerWithOptions_Basic(t *testing.T) {
	t.Parallel()

	logger, err := InitLoggerWithOptions(LoggerOptions{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if logger == nil {
		t.Fatalf("expected non-nil logger")
	}
}

func TestInitLoggerWithOptions_InvalidLevel(t *testing.T) {
	t.Parallel()

	logger, err := InitLoggerWithOptions(LoggerOptions{
		ServiceName: "test-service",
		Environment: "dev",
		Level:       "not-a-level", // triggers the fallback path
		Format:      "json",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if logger == nil {
		t.Fatalf("expected non-nil logger")
	}
}

func TestInitLoggerWithOptions_WithFileAndSampling(t *testing.T) {
	t.Parallel()

	logger, err := InitLoggerWithOptions(LoggerOptions{
		ServiceName: "test-service",
		Environment: "prod", // used together with Sample to hit sampler path
		Level:       "debug",
		Format:      "console",         // when env is prod, this still results in JSON encoder
		FilePath:    "test-logger.log", // non-empty to hit lumberjack configuration
		MaxSize:     1,
		MaxBackups:  2,
		MaxAge:      1,
		Compress:    true,
		Sample:      true,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if logger == nil {
		t.Fatalf("expected non-nil logger")
	}
}

func TestInitLoggerBuffer_VerifiesBufferSink(t *testing.T) {
	t.Parallel()

	logger, buf, err := InitLoggerBuffer("buffered")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if logger == nil {
		t.Fatalf("expected non-nil logger")
	}
	if buf == nil {
		t.Fatalf("expected non-nil buffer")
	}

	logger.Info("hello-buffer", zap.String("k", "v"))
	_ = logger.Sync()

	out := buf.String()
	if out == "" {
		t.Fatalf("expected buffer to contain log output")
	}
	if !containsAll(out, []string{"hello-buffer", "\"k\"", "\"v\""}) {
		t.Fatalf("expected log output to contain message and fields, got: %s", out)
	}
}

func containsAll(s string, subs []string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

func TestAPILoggerMiddleware_RespectsExistingRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	cfg := gateway.LoggingConfig{
		Level:       "info",
		Format:      "json",
		Environment: "dev",
	}

	existingID := uuid.New().String()

	r.Use(APILoggerMiddleware(cfg, "gateway-test"))
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("X-Request-Id", existingID)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if got := w.Header().Get("X-Request-Id"); got != "" {
		t.Fatalf("expected middleware to leave X-Request-Id response header untouched, got %q", got)
	}
}
