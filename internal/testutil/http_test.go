package testutil

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGinEngine_ReturnsEngineInTestMode(t *testing.T) {
	t.Parallel()

	engine := NewGinEngine()
	require.NotNil(t, engine)
}

func TestNewJSONRequest_WithBody(t *testing.T) {
	t.Parallel()

	type payload struct {
		Name string `json:"name"`
	}

	req := NewJSONRequest(t, http.MethodPost, "/test", payload{Name: "doc"})

	require.Equal(t, http.MethodPost, req.Method)
	require.Equal(t, "/test", req.URL.Path)
	require.Equal(t, "application/json", req.Header.Get("Content-Type"))
}

func TestNewJSONRequest_NilBody(t *testing.T) {
	t.Parallel()

	req := NewJSONRequest(t, http.MethodGet, "/test", nil)

	require.Equal(t, http.MethodGet, req.Method)
	require.Equal(t, "/test", req.URL.Path)
	require.Equal(t, "application/json", req.Header.Get("Content-Type"))
}
