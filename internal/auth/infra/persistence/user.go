package persistence

import (
	"context"

	"github.com/a1y/doc-formatter/internal/auth/domain/entity"
	"github.com/a1y/doc-formatter/internal/auth/domain/repository"
	"gorm.io/gorm"
)

var _ repository.UserRepository = &userRepository{}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var model UserModel
	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}
	return model.ToEntity()
}

func (r *userRepository) Create(ctx context.Context, dataEntity *entity.User) error {
	err := dataEntity.Validate()
	if err != nil {
		return err
	}
	var dataModel UserModel
	if err := dataModel.FromEntity(dataEntity); err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.WithContext(ctx).Create(&dataModel).Error
		if err != nil {
			return err
		}

		dataEntity.ID = dataModel.ID

		return nil
	})
}
