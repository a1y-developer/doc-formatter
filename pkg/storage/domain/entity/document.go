package entity

import (
	"errors"

	"github.com/google/uuid"
)

type Document struct {
	ID        uuid.UUID `yaml:"id" json:"id"`
	UserID    uuid.UUID `yaml:"userID" json:"userID"`
	FileName  string    `yaml:"fileName" json:"fileName"`
	FileSize  int64     `yaml:"fileSize" json:"fileSize"`
	ObjectKey string    `yaml:"objectKey" json:"objectKey"`
}

func (d *Document) Validate() error {
	if d.UserID == uuid.Nil {
		return errors.New("user id is required")
	}
	if d.FileName == "" {
		return errors.New("file name is required")
	}
	if d.ObjectKey == "" {
		return errors.New("object key is required")
	}
	return nil
}
