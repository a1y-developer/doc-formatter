package repository

import (
	"context"

	"github.com/a1y/doc-formatter/internal/storage/domain/entity"
	"github.com/google/uuid"
)

type DocumentRepository interface {
	Create(ctx context.Context, d *entity.Document) error
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Document, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Document, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
