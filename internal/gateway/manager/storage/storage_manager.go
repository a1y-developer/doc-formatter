package storage

import (
	"context"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/a1y/doc-formatter/internal/gateway/domain/response"
)

func (m *StorageManager) UploadFile(ctx context.Context, userID string, fileName string, fileSize int64, content []byte) (*response.UploadFileResponse, error) {
	req := &storagepb.UploadFileRequest{
		UserId:   userID,
		FileName: fileName,
		FileSize: fileSize,
		Content:  content,
	}
	resp, err := m.client.UploadFile(ctx, req)
	if err != nil {
		return nil, err
	}
	return &response.UploadFileResponse{
		FileID:   resp.GetFileId(),
		FileName: resp.GetFileName(),
	}, nil
}
