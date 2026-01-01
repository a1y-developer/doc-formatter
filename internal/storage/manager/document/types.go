package document

import (
	"github.com/a1y/doc-formatter/internal/storage/domain/repository"
	"github.com/a1y/doc-formatter/internal/storage/util/s3"
)

type DocumentManager struct {
	documentRepo repository.DocumentRepository
	s3Storage    *s3.S3Storage
}

func NewDocumentManager(
	documentRepo repository.DocumentRepository,
	s3Storage *s3.S3Storage,
) *DocumentManager {
	return &DocumentManager{
		documentRepo: documentRepo,
		s3Storage:    s3Storage,
	}
}
