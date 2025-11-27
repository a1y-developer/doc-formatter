package persistence

import (
	"github.com/a1y/doc-formatter/internal/auth/domain/entity"
)

type UserModel struct {
	BaseModel
	Name     string
	Username string `gorm:"index:unique_user,unique"`
	Email    string `gorm:"index:unique_user,unique"`
	Password string `json:"-"`
}

func (u *UserModel) TableName() string {
	return "users"
}

func (u *UserModel) ToEntity() (*entity.User, error) {
	return &entity.User{
		ID:    u.ID,
		Email: u.Email,
	}, nil
}

func (u *UserModel) FromEntity(e *entity.User) error {
	u.ID = e.ID
	u.Email = e.Email
	return nil
}
