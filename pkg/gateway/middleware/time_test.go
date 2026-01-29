package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTime_SetsStartTimeInContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	var capturedTime time.Time
	Time(c)

	capturedTime = GetStartTime(c.Request.Context())

	assert.False(t, capturedTime.IsZero(), "start time should be set")
	assert.True(t, time.Since(capturedTime) < time.Second, "start time should be recent")
}

func TestTime_DoesNotOverwriteExistingStartTime(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

	// Set a start time manually
	pastTime := time.Now().Add(-5 * time.Second)
	ctx := context.WithValue(c.Request.Context(), StartTimeKey, pastTime)
	c.Request = c.Request.WithContext(ctx)

	Time(c)

	capturedTime := GetStartTime(c.Request.Context())
	assert.Equal(t, pastTime.Unix(), capturedTime.Unix(), "start time should not be overwritten")
}

func TestGetStartTime_WithNilContext(t *testing.T) {
	result := GetStartTime(context.Background())
	assert.True(t, result.IsZero(), "should return zero time for nil context")
}

func TestGetStartTime_WithNoStartTime(t *testing.T) {
	ctx := context.Background()
	result := GetStartTime(ctx)
	assert.True(t, result.IsZero(), "should return zero time when start time not set")
}

func TestGetStartTime_WithStartTime(t *testing.T) {
	expectedTime := time.Now()
	ctx := context.WithValue(context.Background(), StartTimeKey, expectedTime)

	result := GetStartTime(ctx)
	assert.Equal(t, expectedTime.Unix(), result.Unix(), "should return the correct start time")
}

func TestGetStartTime_WithWrongType(t *testing.T) {
	ctx := context.WithValue(context.Background(), StartTimeKey, "not a time")

	result := GetStartTime(ctx)
	assert.True(t, result.IsZero(), "should return zero time when value is wrong type")
}

func TestTime_CallsNext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, router := gin.CreateTestContext(w)

	nextCalled := false
	router.Use(Time)
	router.GET("/test", func(c *gin.Context) {
		nextCalled = true
		startTime := GetStartTime(c.Request.Context())
		assert.False(t, startTime.IsZero())
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, c.Request)

	assert.True(t, nextCalled, "handler should be called after Time middleware")
}
