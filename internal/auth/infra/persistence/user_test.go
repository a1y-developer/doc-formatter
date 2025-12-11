package persistence_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a1y/doc-formatter/internal/auth/domain/entity"
	"github.com/a1y/doc-formatter/internal/auth/infra"
	"github.com/a1y/doc-formatter/internal/auth/infra/persistence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := infra.GetMockDB()
	assert.NoError(t, err)

	repo := persistence.NewUserRepository(db)
	ctx := context.Background()

	user := &entity.User{
		ID:         uuid.New(),
		Email:      "test@example.com",
		Password:   "hashedpassword",
		IsVerified: false,
	}

	// GORM uses Query with RETURNING clause for INSERT
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.Email, user.Password, user.IsVerified).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))
	mock.ExpectClose()

	err = repo.Create(ctx, user)
	assert.NoError(t, err)

	infra.CloseDB(t, db)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, err := infra.GetMockDB()
	assert.NoError(t, err)

	repo := persistence.NewUserRepository(db)
	ctx := context.Background()

	email := "existing@example.com"
	user := &entity.User{
		ID:         uuid.New(),
		Email:      email,
		Password:   "hashedpassword",
		IsVerified: false,
	}

	// GORM will select all fields from UserModel including BaseModel fields
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "description", "name", "username", "email", "password", "is_verified"}).
		AddRow(user.ID.String(), nil, nil, nil, "", "", "", user.Email, user.Password, user.IsVerified)

	// GORM adds soft delete check (deleted_at IS NULL)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)

	foundUser, err := repo.GetByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.IsVerified, foundUser.IsVerified)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("nonexistent@example.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err = repo.GetByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mock.ExpectClose()
	infra.CloseDB(t, db)

	assert.NoError(t, mock.ExpectationsWereMet())
}
