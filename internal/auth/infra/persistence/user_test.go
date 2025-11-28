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
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(user.ID, user.Email, user.Password, user.IsVerified).
		WillReturnResult(sqlmock.NewResult(1, 1))
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
		ID:       uuid.New(),
		Email:    email,
		Password: "hashedpassword",
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password", "is_verified"}).
		AddRow(user.ID, user.Email, user.Password, user.IsVerified)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)

	foundUser, err := repo.GetByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.Email, foundUser.Email)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("nonexistent@example.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err = repo.GetByEmail(ctx, "nonexistent@example.com")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mock.ExpectClose()
	infra.CloseDB(t, db)

	assert.NoError(t, mock.ExpectationsWereMet())
}
