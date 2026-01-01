package route

import (
	"testing"

	"github.com/a1y/doc-formatter/internal/gateway"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := &gateway.Config{
		Address:        ":8080",
		AuthService:    ":8081",
		StorageService: ":8082",
	}

	r, err := NewRouter(config)
	assert.NoError(t, err)
	assert.NotNil(t, r)

	routes := r.Routes()
	expectedRoutes := map[string]string{
		"/api/v1/auth/signup":    "POST",
		"/api/v1/auth/login":     "POST",
		"/api/v1/storage/upload": "POST",
		"/swagger/*any":          "GET",
	}

	for _, route := range routes {
		if method, ok := expectedRoutes[route.Path]; ok {
			assert.Equal(t, method, route.Method)
			delete(expectedRoutes, route.Path)
		}
	}

	assert.Empty(t, expectedRoutes, "Some expected routes were not found")
}
