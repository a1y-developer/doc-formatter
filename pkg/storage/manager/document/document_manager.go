package document

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/a1y/doc-formatter/pkg/storage/domain/entity"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

func (m *DocumentManager) UploadDocument(ctx context.Context, document *entity.Document, file io.Reader) (*entity.Document, error) {
	var createdEntity entity.Document
	if err := copier.Copy(&createdEntity, &document); err != nil {
		return nil, err
	}

	createdEntity.ObjectKey = fmt.Sprintf("%s/%s", createdEntity.OwnerID.String(), createdEntity.Name)

	ok, err := m.s3Storage.PutObject(ctx, createdEntity.ObjectKey, file)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("failed to upload document")
	}

	if err := m.documentRepo.Create(ctx, &createdEntity); err != nil {
		return nil, err
	}
	return &createdEntity, nil
}

func (m *DocumentManager) CreateDocument(ctx context.Context, document *entity.Document) (*entity.Document, error) {
	if err := m.documentRepo.Create(ctx, document); err != nil {
		return nil, err
	}
	return document, nil
}

func (m *DocumentManager) ListDocumentsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Document, error) {
	return m.documentRepo.ListByUserID(ctx, userID)
}

func (m *DocumentManager) GetDocumentByUserID(ctx context.Context, userID uuid.UUID, fileID uuid.UUID) (*entity.Document, error) {
	return m.documentRepo.GetByUserID(ctx, userID, fileID)
}
