package storage

import (
	"context"
	"testing"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/stretchr/testify/require"
)

type fakeStorageClient struct{}

// ListFilesByUserId implements [storage.StorageClient].
func (f *fakeStorageClient) ListFilesByUserId(ctx context.Context, req *storagepb.ListFilesByUserIdRequest) (*storagepb.ListFilesByUserIdResponse, error) {
	return &storagepb.ListFilesByUserIdResponse{
		Documents: []*storagepb.Document{
			{
				Id:       "file-id-1",
				UserId:   req.UserId,
				FileName: "file1.txt",
				FileSize: 100,
			},
			{
				Id:       "file-id-2",
				UserId:   req.UserId,
				FileName: "file2.txt",
				FileSize: 200,
			},
		},
	}, nil
}

func (f *fakeStorageClient) UploadFile(ctx context.Context, req *storagepb.UploadFileRequest) (*storagepb.UploadFileResponse, error) {
	return &storagepb.UploadFileResponse{
		Document: &storagepb.Document{
			Id:       "fake-file-id",
			UserId:   req.UserId,
			FileName: req.FileName,
			FileSize: req.FileSize,
		},
	}, nil
}

func TestNewStorageManager_ReturnsManagerWithClient(t *testing.T) {
	t.Parallel()

	client := &fakeStorageClient{}
	mgr := NewStorageManager(client)

	require.NotNil(t, mgr)
	require.Equal(t, client, mgr.client)
}
