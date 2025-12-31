package persistence

import (
	"gorm.io/gorm"
)

// AutoMigrate runs database migrations for the storage service.
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&DocumentModel{}); err != nil {
		return err
	}
	return nil
}
