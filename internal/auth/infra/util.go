package infra

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a1y/doc-formatter/internal/auth/infra/persistence"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&persistence.UserModel{}); err != nil {
		return err
	}
	return nil
}

// MultiString is a custom type for handling arrays of strings with GORM.
type MultiString []string

// Scan implements the Scanner interface for the MultiString type.
func (s *MultiString) Scan(src any) error {
	switch src := src.(type) {
	case []byte:
		*s = strings.Split(string(src), ",")
	case string:
		*s = strings.Split(src, ",")
	case nil:
		*s = nil
	default:
		return fmt.Errorf("unsupported type %T", src)
	}
	return nil
}

// Value implements the Valuer interface for the MultiString type.
func (s MultiString) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return strings.Join(s, ","), nil
}

// GormDataType gorm common data type
func (s MultiString) GormDataType() string {
	return "text"
}

// GormDBDataType gorm db data type
func (s MultiString) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	// returns different database type based on driver name
	switch db.Name() {
	case "postgres", "sqlite":
		return "text"
	}
	return ""
}

// Create a mock database connection
func GetMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	// Create a sqlMock of sql.DB.
	fakeDB, sqlMock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	// common execution for orm

	// Create the gorm database connection with fake db
	fakeGDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: fakeDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, nil, err
	}

	return fakeGDB, sqlMock, nil
}

// Close the gorm database connection
func CloseDB(t *testing.T, gdb *gorm.DB) {
	db, err := gdb.DB()
	require.NoError(t, err)
	require.NoError(t, db.Close())
}
