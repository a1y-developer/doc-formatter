package handler

import (
	"bytes"
	"context"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/a1y/doc-formatter/pkg/storage/domain/entity"
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
		Document: &storagepb.Document{
			Id:        documentResponse.ID.String(),
			UserId:    documentResponse.UserID.String(),
			FileName:  documentResponse.FileName,
			FileSize:  documentResponse.FileSize,
			ObjectKey: documentResponse.ObjectKey,
			CreatedAt: documentResponse.CreateAt.Unix(),
			UpdatedAt: documentResponse.UpdateAt.Unix(),
		},
	}, nil
}

func (h *Handler) ListFilesByUserId(ctx context.Context, req *storagepb.ListFilesByUserIdRequest) (*storagepb.ListFilesByUserIdResponse, error) {
	userID := uuid.MustParse(req.UserId)
	documents, err := h.documentManager.ListDocumentsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	pbDocuments := make([]*storagepb.Document, len(documents))
	for i, doc := range documents {
		pbDocuments[i] = &storagepb.Document{
			Id:        doc.ID.String(),
			UserId:    doc.UserID.String(),
			FileName:  doc.FileName,
			FileSize:  doc.FileSize,
			ObjectKey: doc.ObjectKey,
			CreatedAt: doc.CreateAt.Unix(),
			UpdatedAt: doc.UpdateAt.Unix(),
		}
	}

	return &storagepb.ListFilesByUserIdResponse{
		Documents: pbDocuments,
	}, nil
}
