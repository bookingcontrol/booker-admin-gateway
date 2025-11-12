package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	uc "github.com/bookingcontrol/booker-admin-gateway/internal/usecase/auth"
)

// MockAuthRepository is a mock for auth repository (used in use case)
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

func TestAuthHandler_Register(t *testing.T) {
	e := echo.New()

	t.Run("successful registration", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		svc := uc.NewService(mockRepo)
		handler := NewAuthHandler(svc)

		reqBody := map[string]interface{}{
			"username": "testuser",
			"password": "password123",
			"email":    "test@example.com",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("UserExists", mock.Anything, "testuser").Return(false, nil)
		mockRepo.On("CreateUser", mock.Anything, "testuser", mock.AnythingOfType("map[string]interface {}")).Return(nil)

		err := handler.Register(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response uc.RegisterView
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "testuser", response.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		svc := uc.NewService(mockRepo)
		handler := NewAuthHandler(svc)

		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte("invalid json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.Register(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockRepo.AssertNotCalled(t, "UserExists")
	})

	t.Run("username already exists", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		svc := uc.NewService(mockRepo)
		handler := NewAuthHandler(svc)

		reqBody := map[string]interface{}{
			"username": "existinguser",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("UserExists", mock.Anything, "existinguser").Return(true, nil)

		err := handler.Register(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "CreateUser")
	})
}

func TestAuthHandler_Login(t *testing.T) {
	e := echo.New()

	t.Run("successful login", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		svc := uc.NewService(mockRepo)
		handler := NewAuthHandler(svc)

		reqBody := map[string]interface{}{
			"username": "testuser",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("UserExists", mock.Anything, "testuser").Return(true, nil)
		mockRepo.On("GetUserPassword", mock.Anything, "testuser").Return("password123", nil)

		err := handler.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		var response uc.LoginView
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, "token-testuser", response.AccessToken)
		assert.Equal(t, "refresh-testuser", response.RefreshToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		svc := uc.NewService(mockRepo)
		handler := NewAuthHandler(svc)

		reqBody := map[string]interface{}{
			"username": "testuser",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("UserExists", mock.Anything, "testuser").Return(true, nil)
		mockRepo.On("GetUserPassword", mock.Anything, "testuser").Return("correctpassword", nil)

		err := handler.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(MockAuthRepository)
		svc := uc.NewService(mockRepo)
		handler := NewAuthHandler(svc)

		reqBody := map[string]interface{}{
			"username": "nonexistent",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("UserExists", mock.Anything, "nonexistent").Return(false, nil)

		err := handler.Login(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "GetUserPassword")
	})
}

