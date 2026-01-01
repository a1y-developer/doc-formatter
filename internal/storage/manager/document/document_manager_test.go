package document

import (
	"bytes"
	"context"
	"testing"

	"github.com/a1y/doc-formatter/internal/storage/domain/entity"
	"github.com/a1y/doc-formatter/internal/storage/domain/repository"
	s3util "github.com/a1y/doc-formatter/internal/storage/util/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type mockDocumentRepository struct{}

func (m *mockDocumentRepository) Create(ctx context.Context, d *entity.Document) error {
	return nil
}

func (m *mockDocumentRepository) ListByUserID(ctx context.Context, id uuid.UUID) ([]*entity.Document, error) {
	return nil, nil
}

func (m *mockDocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	return nil, nil
}

func (m *mockDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

var _ repository.DocumentRepository = (*mockDocumentRepository)(nil)

func TestNewDocumentManager(t *testing.T) {
	t.Parallel()

	manager := NewDocumentManager(&mockDocumentRepository{}, &s3util.S3Storage{})
	require.NotNil(t, manager)
}

func TestDocumentManager_UploadDocument_PanicsWithNilS3Storage(t *testing.T) {
	t.Parallel()

	doc := &entity.Document{
		ID:       uuid.New(),
		UserID:   uuid.New(),
		FileName: "file.txt",
		FileSize: 10,
	}

	manager := &DocumentManager{
		documentRepo: &mockDocumentRepository{},
	}

	reader := bytes.NewReader([]byte("content"))

	require.Panics(t, func() {
		_, _ = manager.UploadDocument(context.Background(), doc, reader)
	})
}
