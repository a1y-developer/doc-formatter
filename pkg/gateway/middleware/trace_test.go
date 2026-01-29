package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTraceID_WithExistingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set(TraceIDHeader, "custom-trace-id-123")

	TraceID(c)

	traceID := GetTraceID(c.Request.Context())
	assert.Equal(t, "custom-trace-id-123", traceID, "should use trace ID from header")
}

func TestTraceID_WithoutHeader_GeneratesTraceID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	TraceID(c)

	traceID := GetTraceID(c.Request.Context())
	assert.NotEmpty(t, traceID, "should generate a trace ID")
	assert.Contains(t, traceID, "-", "generated trace ID should contain hyphen")
}

func TestTraceID_GeneratesUniqueIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var traceIDs []string
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		TraceID(c)

		traceID := GetTraceID(c.Request.Context())
		traceIDs = append(traceIDs, traceID)
	}

	uniqueIDs := make(map[string]bool)
	for _, id := range traceIDs {
		uniqueIDs[id] = true
	}
	assert.Equal(t, 5, len(uniqueIDs), "should generate unique trace IDs")
}

func TestTraceID_GeneratedIDFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	TraceID(c)

	traceID := GetTraceID(c.Request.Context())
	parts := strings.Split(traceID, "-")

	assert.GreaterOrEqual(t, len(parts), 2, "trace ID should have at least 2 parts")

	lastPart := parts[len(parts)-1]
	assert.Len(t, lastPart, 6, "counter part should be 6 digits")
}

func TestGetTraceID_WithNilContext(t *testing.T) {
	result := GetTraceID(context.TODO())
	assert.Equal(t, "", result, "should return empty string for nil context")
}

func TestGetTraceID_WithNoTraceID(t *testing.T) {
	ctx := context.Background()
	result := GetTraceID(ctx)
	assert.Equal(t, "", result, "should return empty string when trace ID not set")
}

func TestGetTraceID_WithTraceID(t *testing.T) {
	expectedID := "test-trace-123"
	ctx := context.WithValue(context.Background(), TraceIDKey, expectedID)

	result := GetTraceID(ctx)
	assert.Equal(t, expectedID, result, "should return the correct trace ID")
}

func TestGetTraceID_WithWrongType(t *testing.T) {
	ctx := context.WithValue(context.Background(), TraceIDKey, 12345)

	result := GetTraceID(ctx)
	assert.Equal(t, "", result, "should return empty string when value is wrong type")
}

func TestTraceID_CallsNext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	nextCalled := false
	router.Use(TraceID)
	router.GET("/test", func(c *gin.Context) {
		nextCalled = true
		traceID := GetTraceID(c.Request.Context())
		assert.NotEmpty(t, traceID)
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, c.Request)

	assert.True(t, nextCalled, "handler should be called after TraceID middleware")
}

func TestTraceIDHeader_IsConfigurable(t *testing.T) {
	assert.Equal(t, "x-df-trace", TraceIDHeader, "TraceIDHeader should be x-df-trace")
}

func TestTraceID_EmptyHeaderGeneratesID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Header.Set(TraceIDHeader, "")

	TraceID(c)

	traceID := GetTraceID(c.Request.Context())
	assert.NotEmpty(t, traceID, "should generate trace ID when header is empty")
}

func TestPrefix_IsInitialized(t *testing.T) {
	assert.NotEmpty(t, prefix, "prefix should be initialized")
	assert.Contains(t, prefix, "/", "prefix should contain slash separator")
}
