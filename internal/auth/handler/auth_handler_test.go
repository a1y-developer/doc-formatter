package handler

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	authpb "github.com/a1y/doc-formatter/api/grpc/auth/v1"
	"github.com/a1y/doc-formatter/internal/auth/infra"
	"github.com/a1y/doc-formatter/internal/auth/infra/persistence"
	"github.com/a1y/doc-formatter/internal/auth/manager/user"
	"github.com/a1y/doc-formatter/pkg/credentials"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Signup(t *testing.T) {
	db, mock, err := infra.GetMockDB()
	assert.NoError(t, err)

	userRepo := persistence.NewUserRepository(db)
	userManager := user.NewUserManager(userRepo)
	h, err := NewHandler(userManager)
	assert.NoError(t, err)

	ctx := context.Background()
	req := &authpb.SignupRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Expect INSERT
	// We use AnyArg() for ID and Password because they are generated/hashed
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(sqlmock.AnyArg(), req.Email, sqlmock.AnyArg(), false).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectClose()

	resp, err := h.Signup(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.UserId)

	infra.CloseDB(t, db)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHandler_Login(t *testing.T) {
	db, mock, err := infra.GetMockDB()
	assert.NoError(t, err)

	userRepo := persistence.NewUserRepository(db)
	userManager := user.NewUserManager(userRepo)
	h, err := NewHandler(userManager)
	assert.NoError(t, err)

	ctx := context.Background()
	password := "password123"

	// Pre-hash password
	hasher := credentials.NewDefaultArgon2idHash()
	hashedPassword, err := hasher.HashPassword(password, nil)
	assert.NoError(t, err)

	userID := uuid.New()
	email := "test@example.com"

	req := &authpb.LoginRequest{
		Email:    email,
		Password: password,
	}

	// Expect SELECT
	rows := sqlmock.NewRows([]string{"id", "email", "password", "is_verified"}).
		AddRow(userID, email, hashedPassword, true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)
	mock.ExpectClose()

	resp, err := h.Login(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotZero(t, resp.ExpiryUnix)

	infra.CloseDB(t, db)
	assert.NoError(t, mock.ExpectationsWereMet())
}
