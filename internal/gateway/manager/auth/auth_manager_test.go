package auth

import (
	"context"
	"testing"

	"github.com/a1y/doc-formatter/internal/gateway/clients/auth"
	"github.com/a1y/doc-formatter/internal/gateway/domain/request"
	"github.com/a1y/doc-formatter/internal/gateway/domain/response"
	"github.com/stretchr/testify/assert"
)

type mockAuthClient struct {
	signupFunc func(ctx context.Context, email, password string) (*response.SignUpResponse, error)
	loginFunc  func(ctx context.Context, email, password string) (*response.LoginResponse, error)
}

func (m *mockAuthClient) Signup(ctx context.Context, email, password string) (*response.SignUpResponse, error) {
	return m.signupFunc(ctx, email, password)
}

func (m *mockAuthClient) Login(ctx context.Context, email, password string) (*response.LoginResponse, error) {
	return m.loginFunc(ctx, email, password)
}

var _ auth.AuthClient = (*mockAuthClient)(nil)

func TestAuthManager_Signup_DelegatesToClient(t *testing.T) {
	t.Parallel()

	expected := &response.SignUpResponse{UserID: "user-123"}
	mockClient := &mockAuthClient{
		signupFunc: func(ctx context.Context, email, password string) (*response.SignUpResponse, error) {
			assert.Equal(t, "test@example.com", email)
			assert.Equal(t, "password", password)
			return expected, nil
		},
		loginFunc: func(ctx context.Context, email, password string) (*response.LoginResponse, error) {
			t.Fatalf("unexpected call to Login")
			return nil, nil
		},
	}

	manager := NewAuthManager(mockClient)

	resp, err := manager.Signup(context.Background(), request.SignupRequest{
		Email:    "test@example.com",
		Password: "password",
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}

func TestAuthManager_Login_DelegatesToClient(t *testing.T) {
	t.Parallel()

	expected := &response.LoginResponse{
		AccessToken: "token-123",
		ExpiryUnix:  42,
	}
	mockClient := &mockAuthClient{
		signupFunc: func(ctx context.Context, email, password string) (*response.SignUpResponse, error) {
			t.Fatalf("unexpected call to Signup")
			return nil, nil
		},
		loginFunc: func(ctx context.Context, email, password string) (*response.LoginResponse, error) {
			assert.Equal(t, "test@example.com", email)
			assert.Equal(t, "password", password)
			return expected, nil
		},
	}

	manager := NewAuthManager(mockClient)

	resp, err := manager.Login(context.Background(), request.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	})

	assert.NoError(t, err)
	assert.Equal(t, expected, resp)
}
