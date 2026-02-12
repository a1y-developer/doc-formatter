package persistence

import (
	"context"

	"github.com/a1y/doc-formatter/pkg/storage/domain/entity"
	"github.com/a1y/doc-formatter/pkg/storage/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ repository.FolderRepository = (*folderRepository)()

type folderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) repository.FolderRepository {
	return &folderRepository{
		db: db,
	}
}

func (r *folderRepository) Create(ctx context.Context, dataEntity *entity.Folder) (*entity.Folder, error) {
	err := dataEntity.Validate()
	if err != nil {
		return nil, err
	}
	var dataModel FolderModel
	if err := dataModel.FromEntity(dataEntity); err != nil {
		return nil, err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		// Create new record in the store
		err = tx.WithContext(ctx).Create(&dataModel).Error
		if err != nil {
			return err
		}
		dataEntity.ID = dataModel.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dataEntity, nil
}

func (r *folderRepository) Get(ctx context.Context, id uuid.UUID) (*entity.Folder, error) {
	var model FolderModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.ToEntity(), nil
}
