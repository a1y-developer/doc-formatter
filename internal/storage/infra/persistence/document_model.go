package persistence

import (
	"github.com/a1y/doc-formatter/internal/storage/domain/entity"
	"github.com/google/uuid"
)

type DocumentModel struct {
	BaseModel
	UserID    uuid.UUID
	FileName  string
	FileSize  int64
	ObjectKey string
}

func (d *DocumentModel) TableName() string {
	return "documents"
}

func (d *DocumentModel) ToEntity() (*entity.Document, error) {
	return &entity.Document{
		ID:        d.ID,
		UserID:    d.UserID,
		FileName:  d.FileName,
		FileSize:  d.FileSize,
		ObjectKey: d.ObjectKey,
	}, nil
}

func (d *DocumentModel) FromEntity(e *entity.Document) error {
	d.ID = e.ID
	d.UserID = e.UserID
	d.FileName = e.FileName
	d.FileSize = e.FileSize
	d.ObjectKey = e.ObjectKey
	return nil
}
