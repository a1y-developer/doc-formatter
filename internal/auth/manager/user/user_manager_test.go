package user

import (
	"context"
	"errors"
	"testing"

	"github.com/a1y/doc-formatter/internal/auth/domain/entity"
	"github.com/a1y/doc-formatter/pkg/credentials"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
