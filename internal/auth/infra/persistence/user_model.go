package persistence

import (
	"github.com/a1y/doc-formatter/internal/auth/domain/entity"
)

type UserModel struct {
	BaseModel
	Name       string
	Username   string `gorm:"index:unique_user,unique"`
	Email      string `gorm:"index:unique_user,unique"`
	Password   string `json:"-"`
	IsVerified bool
}

func (u *UserModel) TableName() string {
	return "users"
}

func (u *UserModel) ToEntity() (*entity.User, error) {
	return &entity.User{
		ID:         u.ID,
		Email:      u.Email,
		Password:   u.Password,
		IsVerified: u.IsVerified,
	}, nil
}

func (u *UserModel) FromEntity(e *entity.User) error {
	u.ID = e.ID
	u.Email = e.Email
	u.Password = e.Password
	u.IsVerified = e.IsVerified
	return nil
}
