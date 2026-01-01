package persistence

import (
	"context"

	"github.com/a1y/doc-formatter/internal/storage/domain/entity"
	"github.com/a1y/doc-formatter/internal/storage/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ repository.DocumentRepository = &documentRepository{}

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) repository.DocumentRepository {
	return &documentRepository{
		db: db,
	}
}

func (r *documentRepository) Create(ctx context.Context, dataEntity *entity.Document) error {
	err := dataEntity.Validate()
	if err != nil {
		return err
	}

	var dataModel DocumentModel
	if err := dataModel.FromEntity(dataEntity); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create new record in the store
		err = tx.WithContext(ctx).Create(&dataModel).Error
		if err != nil {
			return err
		}

		dataEntity.ID = dataModel.ID

		return nil
	})
}

func (r *documentRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Document, error) {
	var models []DocumentModel
	if err := r.db.Where("user_id = ?", userID).Find(&models).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	entities := make([]*entity.Document, len(models))
	for i, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}

func (r *documentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Document, error) {
	var model DocumentModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToEntity()
}

func (r *documentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&DocumentModel{}).Error
}
