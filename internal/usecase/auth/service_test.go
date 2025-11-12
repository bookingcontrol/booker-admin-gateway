package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockAuthRepository is a mock implementation of auth repository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) UserExists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) GetUserPassword(ctx context.Context, username string) (string, error) {
	args := m.Called(ctx, username)
	return args.String(0), args.Error(1)
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, username string, userData map[string]interface{}) error {
	args := m.Called(ctx, username, userData)
	return args.Error(0)
}

func TestService_Register(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		mockRepo.On("UserExists", mock.Anything, "testuser").Return(false, nil)
		mockRepo.On("CreateUser", mock.Anything, "testuser", mock.AnythingOfType("map[string]interface {}")).Return(nil)

		view, err := service.Register(context.Background(), CreateInput{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
		})

		require.NoError(t, err)
		assert.Equal(t, "testuser", view.Username)
		assert.Equal(t, "User registered successfully", view.Message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("missing username", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		_, err := service.Register(context.Background(), CreateInput{
			Username: "",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Equal(t, "username and password are required", err.Error())
		mockRepo.AssertNotCalled(t, "UserExists")
	})

	t.Run("missing password", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		_, err := service.Register(context.Background(), CreateInput{
			Username: "testuser",
			Password: "",
		})

		assert.Error(t, err)
		assert.Equal(t, "username and password are required", err.Error())
		mockRepo.AssertNotCalled(t, "UserExists")
	})

	t.Run("username already exists", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		mockRepo.On("UserExists", mock.Anything, "existinguser").Return(true, nil)

		_, err := service.Register(context.Background(), CreateInput{
			Username: "existinguser",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Equal(t, "username already exists", err.Error())
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "CreateUser")
	})

	t.Run("repository error on UserExists", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		mockRepo.On("UserExists", mock.Anything, "testuser").Return(false, errors.New("db error"))

		_, err := service.Register(context.Background(), CreateInput{
			Username: "testuser",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Equal(t, "internal server error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error on CreateUser", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		mockRepo.On("UserExists", mock.Anything, "testuser").Return(false, nil)
		mockRepo.On("CreateUser", mock.Anything, "testuser", mock.AnythingOfType("map[string]interface {}")).Return(errors.New("db error"))

		_, err := service.Register(context.Background(), CreateInput{
			Username: "testuser",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Equal(t, "failed to create user", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestService_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		mockRepo.On("UserExists", mock.Anything, "testuser").Return(true, nil)
		mockRepo.On("GetUserPassword", mock.Anything, "testuser").Return("password123", nil)

		view, err := service.Login(context.Background(), LoginInput{
			Username: "testuser",
			Password: "password123",
		})

		require.NoError(t, err)
		assert.Equal(t, "token-testuser", view.AccessToken)
		assert.Equal(t, "refresh-testuser", view.RefreshToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("missing username", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		_, err := service.Login(context.Background(), LoginInput{
			Username: "",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Equal(t, "username and password are required", err.Error())
		mockRepo.AssertNotCalled(t, "UserExists")
	})

	t.Run("missing password", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		_, err := service.Login(context.Background(), LoginInput{
			Username: "testuser",
			Password: "",
		})

		assert.Error(t, err)
		assert.Equal(t, "username and password are required", err.Error())
		mockRepo.AssertNotCalled(t, "UserExists")
	})

	t.Run("user does not exist", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		mockRepo.On("UserExists", mock.Anything, "nonexistent").Return(false, nil)

		_, err := service.Login(context.Background(), LoginInput{
			Username: "nonexistent",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "GetUserPassword")
	})

	t.Run("wrong password", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		mockRepo.On("UserExists", mock.Anything, "testuser").Return(true, nil)
		mockRepo.On("GetUserPassword", mock.Anything, "testuser").Return("correctpassword", nil)

		_, err := service.Login(context.Background(), LoginInput{
			Username: "testuser",
			Password: "wrongpassword",
		})

		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error on UserExists", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		mockRepo.On("UserExists", mock.Anything, "testuser").Return(false, errors.New("db error"))

		_, err := service.Login(context.Background(), LoginInput{
			Username: "testuser",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Equal(t, "internal server error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error on GetUserPassword", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		mockRepo.On("UserExists", mock.Anything, "testuser").Return(true, nil)
		mockRepo.On("GetUserPassword", mock.Anything, "testuser").Return("", errors.New("db error"))

		_, err := service.Login(context.Background(), LoginInput{
			Username: "testuser",
			Password: "password123",
		})

		assert.Error(t, err)
		assert.Equal(t, "internal server error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestService_RefreshToken(t *testing.T) {
	t.Run("refresh token", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		service := NewService(mockRepo)

		token, err := service.RefreshToken(context.Background(), "refresh-token")

		require.NoError(t, err)
		assert.Equal(t, "dummy-token", token)
	})
}
