package infra

import (
	"github.com/a1y/doc-formatter/internal/auth/infra/persistence"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&persistence.UserModel{}); err != nil {
		return err
	}
	return nil
}
