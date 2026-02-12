package persistence

import (
	"encoding/json"

	"github.com/a1y/doc-formatter/pkg/storage/domain/entity"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type DocumentModel struct {
	BaseModel
	ParentID  *uuid.UUID     `gorm:"type:uuid;index"`
	Parent    *FolderModel   `gorm:"foreignKey:ParentID"`
	OwnerID   uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name      string         `gorm:"not null;uniqueIndex:idx_file_parent_name"`
	Size      int64          `gorm:"not null"`
	MimeType  string         `gorm:"size:100"`
	Status    string         `gorm:"size:20;default:'PENDING';index"`
	ObjectKey string         `gorm:"not null;unique"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
}

func (d *DocumentModel) TableName() string {
	return "documents"
}

func (d *DocumentModel) ToEntity() (*entity.Document, error) {
	return &entity.Document{
		ID:        d.ID,
		ParentID:  *d.ParentID,
		OwnerID:   d.OwnerID,
		Name:      d.Name,
		Size:      d.Size,
		MineType:  d.MimeType,
		Status:    d.Status,
		ObjectKey: d.ObjectKey,
		Metadata:  d.Metadata,
		CreateAt:  d.CreatedAt,
		UpdateAt:  d.UpdatedAt,
	}, nil
}

func (d *DocumentModel) FromEntity(e *entity.Document) error {
	d.ID = e.ID
	d.ParentID = &e.ParentID
	d.OwnerID = e.OwnerID
	d.Name = e.Name
	d.Size = e.Size
	d.ObjectKey = e.ObjectKey
	if e.Metadata != nil {
		bytes, err := json.Marshal(e.Metadata)
		if err != nil {
			return err
		}
		d.Metadata = datatypes.JSON(bytes)
	}
	return nil
}
