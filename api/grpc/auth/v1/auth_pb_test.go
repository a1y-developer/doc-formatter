package authpb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignupRequest_BasicGettersAndString(t *testing.T) {
	t.Parallel()

	req := &SignupRequest{
		Email:    "test@example.com",
		Password: "secret",
	}

	require.Equal(t, "test@example.com", req.GetEmail())
	require.Equal(t, "secret", req.GetPassword())
	// exercise String and ProtoReflect
	require.NotEmpty(t, req.String())
	require.NotNil(t, req.ProtoReflect())
}

func TestLoginRequest_BasicGetters(t *testing.T) {
	t.Parallel()

	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "secret",
	}

	require.Equal(t, "test@example.com", req.GetEmail())
	require.Equal(t, "secret", req.GetPassword())
	require.NotEmpty(t, req.String())
	require.NotNil(t, req.ProtoReflect())
}

func TestLoginResponse_BasicGetters(t *testing.T) {
	t.Parallel()

	resp := &LoginResponse{
		AccessToken: "token",
		ExpiryUnix:  123,
	}

	require.Equal(t, "token", resp.GetAccessToken())
	require.EqualValues(t, 123, resp.GetExpiryUnix())
	require.NotEmpty(t, resp.String())
	require.NotNil(t, resp.ProtoReflect())
}
