package storage

import (
	"context"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/a1y/doc-formatter/pkg/gateway/domain/response"
)

func (m *StorageManager) UploadFile(ctx context.Context, userID string, fileName string, fileSize int64, content []byte) (*response.Document, error) {
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
	return &response.Document{
		ID:        resp.Document.Id,
		UserID:    resp.Document.UserId,
		FileName:  resp.Document.FileName,
		FileSize:  resp.Document.FileSize,
		CreatedAt: resp.Document.CreatedAt,
		UpdatedAt: resp.Document.UpdatedAt,
	}, nil
}

func (m *StorageManager) ListFilesByUserId(ctx context.Context, userID string) ([]response.Document, error) {
	req := &storagepb.ListFilesByUserIdRequest{
		UserId: userID,
	}
	resp, err := m.client.ListFilesByUserId(ctx, req)
	if err != nil {
		return nil, err
	}

	documents := make([]response.Document, len(resp.Documents))
	for i, doc := range resp.Documents {
		documents[i] = response.Document{
			ID:        doc.Id,
			UserID:    doc.UserId,
			FileName:  doc.FileName,
			FileSize:  doc.FileSize,
			CreatedAt: doc.CreatedAt,
			UpdatedAt: doc.UpdatedAt,
		}
	}

	return documents, nil
}
