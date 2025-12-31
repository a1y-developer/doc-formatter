package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func NewGinEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

func NewJSONRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()

	var buf *bytes.Buffer

	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(t, err)
		buf = bytes.NewBuffer(data)
	} else {
		buf = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")

	return req
}
