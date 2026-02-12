package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidFolderOwner = errors.New("folder owner is required")
	ErrInvalidFolderName  = errors.New("folder name is required")
)

type Folder struct {
	ID         uuid.UUID `json:"id" yaml:"id"`
	OwnerID    uuid.UUID `json:"owner_id" yaml:"owner_id"`
	ParentID   uuid.UUID `json:"parent_id,omitempty" yaml:"parent_id,omitempty"`
	Name       string    `json:"name" yaml:"name"`
	PathTokens []string  `json:"path_tokens" yaml:"path_tokens"`
	CreatedAt  time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" yaml:"updated_at"`
}

// Validate Business Rules
func (f *Folder) Validate() error {
	if f.OwnerID == uuid.Nil {
		return ErrInvalidFolderOwner
	}
	if f.Name == "" {
		return ErrInvalidFolderName
	}
	return nil
}
