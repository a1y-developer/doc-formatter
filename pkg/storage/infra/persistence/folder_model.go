package persistence

import (
	"github.com/a1y/doc-formatter/pkg/storage/domain/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type FolderModel struct {
	BaseModel
	OwnerID    uuid.UUID      `gorm:"type:uuid;not null"`
	ParentID   *uuid.UUID     `gorm:"type:uuid;index"`
	Parent     *FolderModel   `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name       string         `gorm:"not null"`
	PathTokens pq.StringArray `gorm:"type:text[]"`
}

func (f *FolderModel) TableName() string {
	return "folders"
}

func (f *FolderModel) ToEntity() *entity.Folder {
	return &entity.Folder{
		ID:         f.ID,
		OwnerID:    f.OwnerID,
		ParentID:   *f.ParentID,
		Name:       f.Name,
		PathTokens: []string(f.PathTokens),
		CreatedAt:  f.CreatedAt,
		UpdatedAt:  f.UpdatedAt,
	}
}

func (f *FolderModel) FromEntity(e *entity.Folder) error {
	f.ID = e.ID
	f.OwnerID = e.OwnerID
	f.ParentID = &e.ParentID
	f.Name = e.Name
	f.PathTokens = pq.StringArray(e.PathTokens)
	return nil
}
