package handler

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
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

// setupTestPrivateKey generates a test RSA private key (PKCS#8 format) and creates a temp file
func setupTestPrivateKey(t *testing.T) *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate test private key: %v", err)
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to marshal PKCS#8 private key: %v", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Create temporary file for private key
	tmpFile, err := os.CreateTemp("", "test-private-key-*.pem")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.Write(privateKeyPEM); err != nil {
		t.Fatalf("Failed to write private key to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Set environment variable to point to the temp file
	os.Setenv("AUTH_JWT_PRIVATE_KEY_PATH", tmpFile.Name())
	
	// Store temp file name for cleanup
	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
		os.Unsetenv("AUTH_JWT_PRIVATE_KEY_PATH")
	})

	return privateKey
}


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

	// Expect INSERT - GORM uses Query with RETURNING clause
	// GORM inserts all fields from UserModel including BaseModel fields
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), req.Email, sqlmock.AnyArg(), false).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.ExpectClose()

	resp, err := h.Signup(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.UserId)

	infra.CloseDB(t, db)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHandler_Login(t *testing.T) {
	setupTestPrivateKey(t)

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

	// Expect SELECT - GORM adds soft delete check (deleted_at IS NULL)
	// GORM selects all fields from UserModel including BaseModel fields
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "description", "name", "username", "email", "password", "is_verified"}).
		AddRow(userID.String(), nil, nil, nil, "", "", "", email, hashedPassword, false)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
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
