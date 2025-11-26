package infra

import (
    "context"

    "github.com/a1y/ai-doc-formatter/internal/auth/domain"
    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
    m := UserModel{
        ID:           u.ID,
        Email:        u.Email,
        Password:     u.Password,
        IsVerified:   u.IsVerified,
    }
    return r.db.WithContext(ctx).Create(&m).Error
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    var m UserModel
    if err := r.db.WithContext(ctx).
        Where("email = ?", email).
        First(&m).Error; err != nil {
        return nil, err
    }

    return &domain.User{
        ID:           m.ID,
        Email:        m.Email,
        Password:     m.Password,
        IsVerified:   m.IsVerified,
    }, nil
}

func (r *UserRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
    var m UserModel
    if err := r.db.WithContext(ctx).
        Where("id = ?", id).
        First(&m).Error; err != nil {
        return nil, err
    }

    return &domain.User{
        ID:           m.ID,
        Email:        m.Email,
        Password:     m.Password,
        IsVerified:   m.IsVerified,
    }, nil
}
