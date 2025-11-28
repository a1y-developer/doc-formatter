package request

import (
	"testing"

	"github.com/a1y/doc-formatter/internal/gateway/domain/constant"
	"github.com/stretchr/testify/assert"
)

func TestSignupRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     SignupRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: SignupRequest{
				Email:    "user@example.com",
				Password: "supersecret",
			},
			wantErr: nil,
		},
		{
			name: "missing email",
			req: SignupRequest{
				Password: "supersecret",
			},
			wantErr: constant.ErrEmptyEmail,
		},
		{
			name: "missing password",
			req: SignupRequest{
				Email: "user@example.com",
			},
			wantErr: constant.ErrEmptyPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestLoginRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     LoginRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: LoginRequest{
				Email:    "user@example.com",
				Password: "supersecret",
			},
			wantErr: nil,
		},
		{
			name: "missing email",
			req: LoginRequest{
				Password: "supersecret",
			},
			wantErr: constant.ErrEmptyEmail,
		},
		{
			name: "missing password",
			req: LoginRequest{
				Email: "user@example.com",
			},
			wantErr: constant.ErrEmptyPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
