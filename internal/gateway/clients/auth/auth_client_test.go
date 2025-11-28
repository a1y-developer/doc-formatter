package auth

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	authpb "github.com/a1y/doc-formatter/api/grpc/auth/v1"
	"github.com/a1y/doc-formatter/internal/gateway/domain/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// MockAuthServiceClient is a mock of authpb.AuthServiceClient
type MockAuthServiceClient struct {
	mock.Mock
}

func (m *MockAuthServiceClient) Signup(ctx context.Context, in *authpb.SignupRequest, opts ...grpc.CallOption) (*authpb.SignupResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authpb.SignupResponse), args.Error(1)
}

func (m *MockAuthServiceClient) Login(ctx context.Context, in *authpb.LoginRequest, opts ...grpc.CallOption) (*authpb.LoginResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authpb.LoginResponse), args.Error(1)
}

func TestAuthClient_Signup(t *testing.T) {
	email := "test@example.com"
	password := "password123"
	userID := "123"

	t.Run("Success", func(t *testing.T) {
		mockClient := new(MockAuthServiceClient)
		client := &authClient{
			client: mockClient,
		}

		mockClient.On("Signup", mock.Anything, &authpb.SignupRequest{
			Email:    email,
			Password: password,
		}, mock.Anything).Return(&authpb.SignupResponse{UserId: userID}, nil)

		resp, err := client.Signup(context.Background(), email, password)

		assert.NoError(t, err)
		assert.Equal(t, &response.SignUpResponse{UserID: userID}, resp)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := new(MockAuthServiceClient)
		client := &authClient{
			client: mockClient,
		}

		expectedErr := errors.New("signup failed")
		mockClient.On("Signup", mock.Anything, &authpb.SignupRequest{
			Email:    email,
			Password: password,
		}, mock.Anything).Return(nil, expectedErr)

		resp, err := client.Signup(context.Background(), email, password)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, expectedErr, err)
		mockClient.AssertExpectations(t)
	})
}

func TestAuthClient_Login(t *testing.T) {
	email := "test@example.com"
	password := "password123"
	accessToken := "access_token"
	expiryUnix := int64(1700000000)

	t.Run("Success", func(t *testing.T) {
		mockClient := new(MockAuthServiceClient)
		client := &authClient{
			client: mockClient,
		}

		mockClient.On("Login", mock.Anything, &authpb.LoginRequest{
			Email:    email,
			Password: password,
		}, mock.Anything).Return(&authpb.LoginResponse{
			AccessToken: accessToken,
			ExpiryUnix:  expiryUnix,
		}, nil)

		resp, err := client.Login(context.Background(), email, password)

		assert.NoError(t, err)
		assert.Equal(t, &response.LoginResponse{
			AccessToken: accessToken,
			ExpiryUnix:  expiryUnix,
		}, resp)
		mockClient.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockClient := new(MockAuthServiceClient)
		client := &authClient{
			client: mockClient,
		}

		expectedErr := errors.New("login failed")
		mockClient.On("Login", mock.Anything, &authpb.LoginRequest{
			Email:    email,
			Password: password,
		}, mock.Anything).Return(nil, expectedErr)

		resp, err := client.Login(context.Background(), email, password)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, expectedErr, err)
		mockClient.AssertExpectations(t)
	})
}

type fakeAuthServiceServer struct {
	authpb.UnimplementedAuthServiceServer
}

func (f *fakeAuthServiceServer) Signup(ctx context.Context, req *authpb.SignupRequest) (*authpb.SignupResponse, error) {
	return &authpb.SignupResponse{UserId: "srv-" + req.GetEmail()}, nil
}

func (f *fakeAuthServiceServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	return &authpb.LoginResponse{
		AccessToken: "token-" + req.GetEmail(),
		ExpiryUnix:  time.Now().Add(time.Minute).Unix(),
	}, nil
}

func TestNewAuthClient(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, &fakeAuthServiceServer{})

	go func() {
		_ = grpcServer.Serve(lis)
	}()
	t.Cleanup(func() {
		grpcServer.Stop()
		_ = lis.Close()
	})

	client := NewAuthClient(lis.Addr().String())
	require.IsType(t, &authClient{}, client)
	t.Cleanup(func() {
		_ = client.(*authClient).conn.Close()
	})

	resp, err := client.Signup(context.Background(), "user@example.com", "secret")
	assert.NoError(t, err)
	assert.Equal(t, "srv-user@example.com", resp.UserID)

	loginResp, err := client.Login(context.Background(), "user@example.com", "secret")
	assert.NoError(t, err)
	assert.Contains(t, loginResp.AccessToken, "token-user@example.com")
	assert.NotZero(t, loginResp.ExpiryUnix)
}
