package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Document struct {
	ID        uuid.UUID      `yaml:"id" json:"id"`
	ParentID  uuid.UUID      `yaml:"parent_id" json:"parent_id"`
	OwnerID   uuid.UUID      `yaml:"ownerID" json:"ownerID"`
	Name      string         `yaml:"name" json:"name"`
	Size      int64          `yaml:"size" json:"size"`
	MineType  string         `yaml:"mineType" json:"mineType"`
	Status    string         `yaml:"status" json:"status"`
	ObjectKey string         `yaml:"objectKey" json:"objectKey"`
	Metadata  datatypes.JSON `yaml:"metadata" json:"metadata"`
	CreateAt  time.Time      `yaml:"createAt" json:"createAt"`
	UpdateAt  time.Time      `yaml:"updateAt" json:"updateAt"`
}

func (d *Document) Validate() error {
	if d.OwnerID == uuid.Nil {
		return errors.New("user id is required")
	}
	if d.Name == "" {
		return errors.New("file name is required")
	}
	if d.ObjectKey == "" {
		return errors.New("object key is required")
	}
	return nil
}
