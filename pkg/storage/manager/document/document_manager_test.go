package document

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/a1y/doc-formatter/pkg/storage/domain/entity"
	"github.com/a1y/doc-formatter/pkg/storage/domain/repository"
	s3util "github.com/a1y/doc-formatter/pkg/storage/util/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type mockDocumentRepository struct {
	documents []*entity.Document
	createErr error
	err       error
}

func (m *mockDocumentRepository) Create(ctx context.Context, d *entity.Document) error {
	if m.createErr != nil {
		return m.createErr
	}
	// Simulate ID assignment
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

func (m *mockDocumentRepository) ListByUserID(ctx context.Context, id uuid.UUID) ([]*entity.Document, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.documents, nil
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

func TestDocumentManager_ListDocumentsByUserID_Success(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	expectedDocs := []*entity.Document{
		{
			ID:       uuid.New(),
			UserID:   userID,
			FileName: "file1.txt",
			FileSize: 100,
		},
		{
			ID:       uuid.New(),
			UserID:   userID,
			FileName: "file2.pdf",
			FileSize: 200,
		},
	}

	mockRepo := &mockDocumentRepository{
		documents: expectedDocs,
	}
	manager := NewDocumentManager(mockRepo, &s3util.S3Storage{})

	ctx := context.Background()
	docs, err := manager.ListDocumentsByUserID(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, docs)
	require.Len(t, docs, 2)
	require.Equal(t, expectedDocs[0].FileName, docs[0].FileName)
	require.Equal(t, expectedDocs[1].FileName, docs[1].FileName)
}

func TestDocumentManager_ListDocumentsByUserID_EmptyResult(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	mockRepo := &mockDocumentRepository{
		documents: []*entity.Document{},
	}
	manager := NewDocumentManager(mockRepo, &s3util.S3Storage{})

	ctx := context.Background()
	docs, err := manager.ListDocumentsByUserID(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, docs)
	require.Empty(t, docs)
}

func TestDocumentManager_ListDocumentsByUserID_PropagatesError(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	expectedErr := errors.New("database error")
	mockRepo := &mockDocumentRepository{
		err: expectedErr,
	}
	manager := NewDocumentManager(mockRepo, &s3util.S3Storage{})

	ctx := context.Background()
	docs, err := manager.ListDocumentsByUserID(ctx, userID)

	require.Error(t, err)
	require.Nil(t, docs)
	require.Equal(t, expectedErr, err)
}

func TestDocumentManager_ListDocumentsByUserID_PanicsWithNilRepository(t *testing.T) {
	t.Parallel()

	manager := &DocumentManager{}
	userID := uuid.New()

	require.Panics(t, func() {
		_, _ = manager.ListDocumentsByUserID(context.Background(), userID)
	})
}
