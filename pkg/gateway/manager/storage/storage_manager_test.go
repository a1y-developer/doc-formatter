package storage

import (
	"context"
	"errors"
	"testing"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/a1y/doc-formatter/pkg/gateway/domain/response"
	"github.com/stretchr/testify/require"
)

type stubStorageClient struct {
	uploadResp *storagepb.UploadFileResponse
	listResp   *storagepb.ListFilesByUserIdResponse
	err        error
}

// ListFilesByUserId implements [storage.StorageClient].
func (s *stubStorageClient) ListFilesByUserId(ctx context.Context, req *storagepb.ListFilesByUserIdRequest) (*storagepb.ListFilesByUserIdResponse, error) {
	return s.listResp, s.err
}

func (s *stubStorageClient) UploadFile(_ context.Context, _ *storagepb.UploadFileRequest) (*storagepb.UploadFileResponse, error) {
	return s.uploadResp, s.err
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
		uploadResp: &storagepb.UploadFileResponse{
			Document: &storagepb.Document{
				Id:       "file-id",
				FileName: "file.txt",
			},
		},
	}
	mgr := NewStorageManager(client)

	ctx := context.Background()
	resp, err := mgr.UploadFile(ctx, "user-id", "file.txt", 123, []byte("content"))

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, &response.Document{
		ID:       "file-id",
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

func TestStorageManager_ListFilesByUserId_Success(t *testing.T) {
	t.Parallel()

	client := &stubStorageClient{
		listResp: &storagepb.ListFilesByUserIdResponse{
			Documents: []*storagepb.Document{
				{
					Id:        "file-1",
					UserId:    "user-id",
					FileName:  "file1.txt",
					FileSize:  100,
					CreatedAt: 1704067200,
					UpdatedAt: 1704067200,
				},
				{
					Id:        "file-2",
					UserId:    "user-id",
					FileName:  "file2.pdf",
					FileSize:  200,
					CreatedAt: 1704153600,
					UpdatedAt: 1704153600,
				},
			},
		},
	}
	mgr := NewStorageManager(client)

	ctx := context.Background()
	resp, err := mgr.ListFilesByUserId(ctx, "user-id")

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp, 2)
	require.Equal(t, "file-1", resp[0].ID)
	require.Equal(t, "file1.txt", resp[0].FileName)
	require.Equal(t, int64(100), resp[0].FileSize)
	require.Equal(t, "file-2", resp[1].ID)
	require.Equal(t, "file2.pdf", resp[1].FileName)
	require.Equal(t, int64(200), resp[1].FileSize)
}

func TestStorageManager_ListFilesByUserId_EmptyResult(t *testing.T) {
	t.Parallel()

	client := &stubStorageClient{
		listResp: &storagepb.ListFilesByUserIdResponse{
			Documents: []*storagepb.Document{},
		},
	}
	mgr := NewStorageManager(client)

	ctx := context.Background()
	resp, err := mgr.ListFilesByUserId(ctx, "user-id")

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Empty(t, resp)
}

func TestStorageManager_ListFilesByUserId_PropagatesError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("list failed")
	client := &stubStorageClient{err: expectedErr}
	mgr := NewStorageManager(client)

	ctx := context.Background()
	resp, err := mgr.ListFilesByUserId(ctx, "user-id")

	require.Error(t, err)
	require.Nil(t, resp)
	require.Equal(t, expectedErr, err)
}
