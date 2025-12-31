package persistence

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a1y/doc-formatter/internal/storage/domain/entity"
	testpersistence "github.com/a1y/doc-formatter/pkg/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDocumentRepository_Create(t *testing.T) {
	db, mock, err := testpersistence.GetMockDB()
	assert.NoError(t, err)

	repo := NewDocumentRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	doc := &entity.Document{
		UserID:    userID,
		FileName:  "test.txt",
		FileSize:  123,
		ObjectKey: userID.String() + "/test.txt",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "documents"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.ExpectCommit()
	mock.ExpectClose()

	err = repo.Create(ctx, doc)
	assert.NoError(t, err)

	testpersistence.CloseDB(t, db)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDocumentRepository_ListByUserID(t *testing.T) {
	db, mock, err := testpersistence.GetMockDB()
	assert.NoError(t, err)

	repo := NewDocumentRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	docID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "description",
		"user_id", "file_name", "file_size", "object_key",
	}).AddRow(
		docID.String(), nil, nil, nil, "",
		userID, "test.txt", int64(123), userID.String()+"/test.txt",
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "documents" WHERE user_id = $1 AND "documents"."deleted_at" IS NULL`)).
		WithArgs(userID).
		WillReturnRows(rows)

	docs, err := repo.ListByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.Equal(t, userID, docs[0].UserID)
	assert.Equal(t, "test.txt", docs[0].FileName)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "documents" WHERE user_id = $1 AND "documents"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at", "description",
			"user_id", "file_name", "file_size", "object_key",
		}))

	emptyDocs, err := repo.ListByUserID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Len(t, emptyDocs, 0)

	mock.ExpectClose()
	testpersistence.CloseDB(t, db)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDocumentRepository_GetByID(t *testing.T) {
	db, mock, err := testpersistence.GetMockDB()
	assert.NoError(t, err)

	repo := NewDocumentRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	docID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at", "description",
		"user_id", "file_name", "file_size", "object_key",
	}).AddRow(
		docID.String(), nil, nil, nil, "",
		userID, "test.txt", int64(123), userID.String()+"/test.txt",
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "documents" WHERE id = $1 AND "documents"."deleted_at" IS NULL ORDER BY "documents"."id" LIMIT $2`)).
		WithArgs(docID, 1).
		WillReturnRows(rows)

	doc, err := repo.GetByID(ctx, docID)
	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, docID, doc.ID)
	assert.Equal(t, userID, doc.UserID)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "documents" WHERE id = $1 AND "documents"."deleted_at" IS NULL ORDER BY "documents"."id" LIMIT $2`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err = repo.GetByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mock.ExpectClose()
	testpersistence.CloseDB(t, db)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDocumentRepository_Delete(t *testing.T) {
	db, mock, err := testpersistence.GetMockDB()
	assert.NoError(t, err)

	repo := NewDocumentRepository(db)
	ctx := context.Background()

	docID := uuid.New()

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "documents" SET "deleted_at"=$1 WHERE id = $2 AND "documents"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), docID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(ctx, docID)
	assert.NoError(t, err)

	mock.ExpectClose()
	testpersistence.CloseDB(t, db)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDocumentRepository_SQLiteIntegration(t *testing.T) {
	t.Parallel()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	require.NoError(t, AutoMigrate(db))

	repo := NewDocumentRepository(db)
	ctx := context.Background()

	userID := uuid.New()
	doc := &entity.Document{
		UserID:    userID,
		FileName:  "test.txt",
		FileSize:  123,
		ObjectKey: userID.String() + "/test.txt",
	}

	// Create
	require.NoError(t, repo.Create(ctx, doc))
	require.NotEqual(t, uuid.Nil, doc.ID)

	// List by user
	docs, err := repo.ListByUserID(ctx, userID)
	require.NoError(t, err)
	require.Len(t, docs, 1)

	// Get by ID
	fetched, err := repo.GetByID(ctx, docs[0].ID)
	require.NoError(t, err)
	require.Equal(t, doc.ID, fetched.ID)
	require.Equal(t, doc.UserID, fetched.UserID)
	require.Equal(t, doc.FileName, fetched.FileName)

	// Delete
	require.NoError(t, repo.Delete(ctx, fetched.ID))
}
