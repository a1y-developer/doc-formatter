package storage

import (
	"context"
	"errors"
	"testing"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/a1y/doc-formatter/internal/gateway/domain/response"
	"github.com/stretchr/testify/require"
)

type stubStorageClient struct {
	resp *storagepb.UploadFileResponse
	err  error
}

func (s *stubStorageClient) UploadFile(_ context.Context, _ *storagepb.UploadFileRequest) (*storagepb.UploadFileResponse, error) {
	return s.resp, s.err
}

func TestNewStorageManager_CreatesManager(t *testing.T) {
	t.Parallel()

	client := &stubStorageClient{}
	mgr := NewStorageManager(client)

	require.NotNil(t, mgr)
	require.Equal(t, client, mgr.client)
}

func TestStorageManager_UploadFile_Success(t *testing.T) {
	t.Parallel()

	client := &stubStorageClient{
		resp: &storagepb.UploadFileResponse{
			FileId:   "file-id",
			FileName: "file.txt",
		},
	}
	mgr := NewStorageManager(client)

	ctx := context.Background()
	resp, err := mgr.UploadFile(ctx, "user-id", "file.txt", 123, []byte("content"))

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, &response.UploadFileResponse{
		FileID:   "file-id",
		FileName: "file.txt",
	}, resp)
}

func TestStorageManager_UploadFile_PropagatesError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("upload failed")
	client := &stubStorageClient{err: expectedErr}
	mgr := NewStorageManager(client)

	ctx := context.Background()
	resp, err := mgr.UploadFile(ctx, "user-id", "file.txt", 123, []byte("content"))

	require.Error(t, err)
	require.Nil(t, resp)
	require.Equal(t, expectedErr, err)
}
