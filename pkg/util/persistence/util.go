package persistence

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetMockDB create a mock database connection
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

// CloseDB close the gorm database connection
func CloseDB(t *testing.T, gdb *gorm.DB) {
	db, err := gdb.DB()
	require.NoError(t, err)
	require.NoError(t, db.Close())
}
