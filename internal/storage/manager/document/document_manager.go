package document

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/a1y/doc-formatter/internal/storage/domain/entity"
	"github.com/jinzhu/copier"
)

func (m *DocumentManager) UploadDocument(ctx context.Context, document *entity.Document, file io.Reader) (*entity.Document, error) {
	var createdEntity entity.Document
	if err := copier.Copy(&createdEntity, &document); err != nil {
		return nil, err
	}

	createdEntity.ObjectKey = fmt.Sprintf("%s/%s", createdEntity.UserID.String(), createdEntity.FileName)

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
