package repository

import (
	"context"

	"github.com/a1y/doc-formatter/pkg/storage/domain/entity"
	"github.com/google/uuid"
)

type DocumentRepository interface {
	Create(ctx context.Context, d *entity.Document) error
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Document, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Document, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, fileID uuid.UUID) (*entity.Document, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type FolderRepository interface {
	Create(ctx context.Context, folder *entity.Folder) (*entity.Folder, error)
	Get(ctx context.Context, id uuid.UUID) (*entity.Folder, error)
}
