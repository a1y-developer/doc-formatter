package handler

import (
	"bytes"
	"context"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/a1y/doc-formatter/internal/storage/domain/entity"
	"github.com/google/uuid"
)

func (h *Handler) UploadFile(ctx context.Context, req *storagepb.UploadFileRequest) (*storagepb.UploadFileResponse, error) {
	reader := bytes.NewReader(req.Content)
	documentEntity := entity.Document{
		UserID:   uuid.MustParse(req.UserId),
		FileName: req.FileName,
		FileSize: req.FileSize,
	}
	documentResponse, err := h.documentManager.UploadDocument(ctx, &documentEntity, reader)
	if err != nil {
		return nil, err
	}
	return &storagepb.UploadFileResponse{
		FileId:   documentResponse.ID.String(),
		FileName: documentResponse.FileName,
	}, nil
}
