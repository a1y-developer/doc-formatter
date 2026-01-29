package request

import (
	"testing"

	"github.com/a1y/doc-formatter/pkg/gateway/domain/constant"
	"github.com/stretchr/testify/assert"
)

func TestUploadFileRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     UploadFileRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: UploadFileRequest{
				UserID: "550e8400-e29b-41d4-a716-446655440000",
				File:   nil, // File validation is done by gin binding, not in Validate()
			},
			wantErr: nil,
		},
		{
			name: "missing user_id",
			req: UploadFileRequest{
				UserID: "",
				File:   nil,
			},
			wantErr: constant.ErrEmptyUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestListFilesByUserIdRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     ListFilesByUserIdRequest
		wantErr error
	}{
		{
			name: "valid request",
			req: ListFilesByUserIdRequest{
				UserID: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: nil,
		},
		{
			name: "missing user_id",
			req: ListFilesByUserIdRequest{
				UserID: "",
			},
			wantErr: constant.ErrEmptyUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
