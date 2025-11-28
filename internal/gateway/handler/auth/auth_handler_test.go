package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a1y/doc-formatter/internal/gateway/domain/request"
	"github.com/a1y/doc-formatter/internal/gateway/domain/response"
	manager "github.com/a1y/doc-formatter/internal/gateway/manager/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthClient is a mock of AuthClient interface
type MockAuthClient struct {
	mock.Mock
}

func (m *MockAuthClient) Signup(ctx context.Context, email, password string) (*response.SignUpResponse, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.SignUpResponse), args.Error(1)
}

func (m *MockAuthClient) Login(ctx context.Context, email, password string) (*response.LoginResponse, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.LoginResponse), args.Error(1)
}

func setupRouter() (*gin.Engine, *MockAuthClient) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockClient := new(MockAuthClient)
	authManager := manager.NewAuthManager(mockClient)
	authHandler, _ := NewAuthHandler(authManager)

	r.POST("/api/auth/signup", authHandler.Signup)
	r.POST("/api/auth/login", authHandler.Login)

	return r, mockClient
}

func TestAuthHandler_Signup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r, mockClient := setupRouter()
		reqBody := request.SignupRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockClient.On("Signup", mock.Anything, reqBody.Email, reqBody.Password).
			Return(&response.SignUpResponse{UserID: "123"}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(jsonBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "123", resp["user_id"])

		mockClient.AssertExpectations(t)
	})

	t.Run("BadRequest_InvalidJSON", func(t *testing.T) {
		r, _ := setupRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/signup", bytes.NewBufferString("invalid-json"))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("BadRequest_ValidationFailed", func(t *testing.T) {
		r, _ := setupRouter()
		reqBody := request.SignupRequest{
			Email:    "invalid-email",
			Password: "123", // too short
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(jsonBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		r, mockClient := setupRouter()
		reqBody := request.SignupRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockClient.On("Signup", mock.Anything, reqBody.Email, reqBody.Password).
			Return(nil, errors.New("internal error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(jsonBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockClient.AssertExpectations(t)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r, mockClient := setupRouter()
		reqBody := request.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockClient.On("Login", mock.Anything, reqBody.Email, reqBody.Password).
			Return(&response.LoginResponse{
				AccessToken: "token",
				ExpiryUnix:  1234567890,
			}, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "token", resp["access_token"])
		assert.Equal(t, float64(1234567890), resp["expiry_unix"])

		mockClient.AssertExpectations(t)
	})

	t.Run("BadRequest_InvalidJSON", func(t *testing.T) {
		r, _ := setupRouter()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString("invalid-json"))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		r, mockClient := setupRouter()
		reqBody := request.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockClient.On("Login", mock.Anything, reqBody.Email, reqBody.Password).
			Return(nil, errors.New("invalid credentials"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockClient.AssertExpectations(t)
	})
}
