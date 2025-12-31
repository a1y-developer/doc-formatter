package storage

import (
	"context"
	"testing"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/stretchr/testify/require"
)

type fakeStorageClient struct{}

func (f *fakeStorageClient) UploadFile(ctx context.Context, req *storagepb.UploadFileRequest) (*storagepb.UploadFileResponse, error) {
	return &storagepb.UploadFileResponse{
		FileId:   "fake-id",
		FileName: req.GetFileName(),
	}, nil
}

func TestNewStorageManager_ReturnsManagerWithClient(t *testing.T) {
	t.Parallel()

	client := &fakeStorageClient{}
	mgr := NewStorageManager(client)

	require.NotNil(t, mgr)
	require.Equal(t, client, mgr.client)
}
