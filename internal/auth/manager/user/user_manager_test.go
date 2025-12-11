package user

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"testing"

	"github.com/a1y/doc-formatter/internal/auth/domain/entity"
	"github.com/a1y/doc-formatter/pkg/credentials"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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


// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *entity.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userManager := NewUserManager(mockRepo)
		user := &entity.User{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
			return u.Email == user.Email && u.Password != "password123" // Password should be hashed
		})).Return(nil)

		createdUser, err := userManager.CreateUser(context.Background(), user)

		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, user.Email, createdUser.Email)
		assert.NotEqual(t, "password123", createdUser.Password)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userManager := NewUserManager(mockRepo)
		user := &entity.User{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

		createdUser, err := userManager.CreateUser(context.Background(), user)

		assert.Error(t, err)
		assert.Nil(t, createdUser)
		mockRepo.AssertExpectations(t)
	})
}

func TestLoginUser(t *testing.T) {
	password := "password123"
	hasher := credentials.NewDefaultArgon2idHash()
	hashedPassword, _ := hasher.HashPassword(password, nil)

	t.Run("Success", func(t *testing.T) {
		setupTestPrivateKey(t)

		mockRepo := new(MockUserRepository)
		userManager := NewUserManager(mockRepo)
		userEntity := &entity.User{
			Email:    "test@example.com",
			Password: password,
		}

		storedUser := &entity.User{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		mockRepo.On("GetByEmail", mock.Anything, userEntity.Email).Return(storedUser, nil)

		token, exp, err := userManager.LoginUser(context.Background(), userEntity)

		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.Greater(t, exp, int64(0))
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		setupTestPrivateKey(t)

		mockRepo := new(MockUserRepository)
		userManager := NewUserManager(mockRepo)
		userEntity := &entity.User{
			Email:    "test@example.com",
			Password: password,
		}

		mockRepo.On("GetByEmail", mock.Anything, userEntity.Email).Return(nil, errors.New("user not found"))

		token, exp, err := userManager.LoginUser(context.Background(), userEntity)

		assert.Error(t, err)
		assert.Nil(t, token)
		assert.Equal(t, int64(0), exp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		setupTestPrivateKey(t)

		mockRepo := new(MockUserRepository)
		userManager := NewUserManager(mockRepo)
		userEntity := &entity.User{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		storedUser := &entity.User{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		mockRepo.On("GetByEmail", mock.Anything, userEntity.Email).Return(storedUser, nil)

		token, exp, err := userManager.LoginUser(context.Background(), userEntity)

		assert.Error(t, err)
		assert.Nil(t, token)
		assert.Equal(t, int64(0), exp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("StoredHashInvalid", func(t *testing.T) {
		setupTestPrivateKey(t)

		mockRepo := new(MockUserRepository)
		userManager := NewUserManager(mockRepo)
		userEntity := &entity.User{
			Email:    "test@example.com",
			Password: password,
		}

		storedUser := &entity.User{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Email:    "test@example.com",
			Password: "not-a-valid-hash",
		}

		mockRepo.On("GetByEmail", mock.Anything, userEntity.Email).Return(storedUser, nil)

		token, exp, err := userManager.LoginUser(context.Background(), userEntity)

		assert.Error(t, err)
		assert.Nil(t, token)
		assert.Equal(t, int64(0), exp)
		mockRepo.AssertExpectations(t)
	})
}
