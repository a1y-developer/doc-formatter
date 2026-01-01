package storage

import (
	"context"
	"time"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
)

func (s *storageClient) UploadFile(ctx context.Context, req *storagepb.UploadFileRequest) (*storagepb.UploadFileResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return s.client.UploadFile(ctx, req)
}
