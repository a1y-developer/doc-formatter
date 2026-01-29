//nolint:dupl
package request

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignupRequestDecode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       string
		wantErr    bool
		wantEmail  string
		wantPasswd string
	}{
		{
			name:       "valid json",
			body:       `{"email":"test@example.com","password":"secret123"}`,
			wantErr:    false,
			wantEmail:  "test@example.com",
			wantPasswd: "secret123",
		},
		{
			name:    "invalid json",
			body:    `{invalid json}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			req := &SignupRequest{}
			err := req.Decode(c)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantEmail, req.Email)
				assert.Equal(t, tt.wantPasswd, req.Password)
			}
		})
	}
}

func TestLoginRequestDecode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       string
		wantErr    bool
		wantEmail  string
		wantPasswd string
	}{
		{
			name:       "valid json",
			body:       `{"email":"user@example.com","password":"password123"}`,
			wantErr:    false,
			wantEmail:  "user@example.com",
			wantPasswd: "password123",
		},
		{
			name:    "invalid json",
			body:    `{not valid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			req := &LoginRequest{}
			err := req.Decode(c)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantEmail, req.Email)
				assert.Equal(t, tt.wantPasswd, req.Password)
			}
		})
	}
}
