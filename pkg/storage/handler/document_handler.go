package handler

import (
	"bytes"
	"context"

	storagepb "github.com/a1y/doc-formatter/api/grpc/storage/v1"
	"github.com/a1y/doc-formatter/pkg/storage/domain/entity"
	"github.com/google/uuid"
)

func (h *Handler) InitiateUploadDocument(ctx context.Context, req *storagepb.InitiateUploadDocumentRequest) (*storagepb.InitiateUploadDocumentResponse, error) {
	// Implementation for initiating upload (e.g., generating pre-signed URL) can be added here.
	docEntity := entity.Document{
		OwnerID:  uuid.MustParse(req.OwnerId),
		Name:     req.Name,
		Size:     req.Size,
		MineType: req.MimeType,
	}
	if req.ParentId == "" {
		folder, err := h.folderManager.CreateFolder(ctx, &entity.Folder{
			OwnerID: uuid.MustParse(req.OwnerId),
		})
		if err != nil {
			return nil, err
		}
		docEntity.ParentID = folder.ID
	} else {
		docEntity.ParentID = uuid.MustParse(req.ParentId)
	}
	document, err := h.documentManager.CreateDocument(ctx, &docEntity)

	return &storagepb.InitiateUploadDocumentResponse{}, nil
}

func (h *Handler) UploadFile(ctx context.Context, req *storagepb.UploadFileRequest) (*storagepb.UploadFileResponse, error) {
	reader := bytes.NewReader(req.Content)
	documentEntity := entity.Document{
		OwnerID:  uuid.MustParse(req.UserId),
		Name:     req.FileName,
		FileSize: req.FileSize,
	}
	documentResponse, err := h.documentManager.UploadDocument(ctx, &documentEntity, reader)
	if err != nil {
		return nil, err
	}
	return &storagepb.UploadFileResponse{
		Document: &storagepb.Document{
			Id:        documentResponse.ID.String(),
			UserId:    documentResponse.OwnerID.String(),
			FileName:  documentResponse.Name,
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
			UserId:    doc.OwnerID.String(),
			FileName:  doc.Name,
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
