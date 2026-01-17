package persistence

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetMockDBAndCloseDB(t *testing.T) {
	t.Parallel()

	db, mock, err := GetMockDB()
	if err != nil {
		t.Fatalf("GetMockDB returned error: %v", err)
	}
	if db == nil {
		t.Fatalf("expected non-nil *gorm.DB from GetMockDB")
	}

	mock.ExpectClose()
	CloseDB(t, db)

	if err := mock.ExpectationsWereMet(); err != nil && err != sqlmock.ErrCancelled {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
