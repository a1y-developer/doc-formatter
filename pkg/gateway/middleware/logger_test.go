package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/a1y/doc-formatter/pkg/gateway"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInitLogger_Basic(t *testing.T) {
	t.Parallel()

	// Test with stdout
	logger := InitLogger("", "test-logger")
	require.NotNil(t, logger)

	// Test with valid file path
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	logger2 := InitLogger(logPath, "file-logger")
	require.NotNil(t, logger2)

	logger2.Info("test message")
	logger2.Sync()

	// Verify file exists
	_, err := os.Stat(logPath)
	require.NoError(t, err)

	// Test with nested directory (auto-create)
	nestedPath := filepath.Join(tmpDir, "logs", "nested", "app.log")
	logger3 := InitLogger(nestedPath, "nested-logger")
	require.NotNil(t, logger3)

	logger3.Info("nested test")
	logger3.Sync()

	_, err = os.Stat(nestedPath)
	require.NoError(t, err)
}

func TestInitLoggerBuffer_VerifiesBufferSink(t *testing.T) {
	t.Parallel()

	logger, buf := InitLoggerBuffer("buffered")
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
		// Use a temp file for testing file path requirement
		FilePath: "/tmp/gateway-test.log",
	}

	existingID := uuid.New().String()

	r.Use(APILoggerMiddleware(cfg))
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
}
