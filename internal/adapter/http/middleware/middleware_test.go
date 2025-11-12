package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/bookingcontrol/booker-admin-gateway/internal/config"
	"github.com/bookingcontrol/booker-admin-gateway/internal/infrastructure/redis"
)

// Используем реальный Redis клиент, но с моком на уровне методов через интерфейс
// Или просто тестируем логику без Redis (только проверка заголовков)

func TestAuthMiddleware(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{}
	// Создаем фиктивный Redis клиент - middleware использует его только для rate limit
	// Для auth middleware Redis не нужен
	redisClient := redis.NewClient("localhost:6379", "")
	defer redisClient.Close()
	
	mw := New(redisClient, cfg)

	t.Run("missing authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := mw.AuthMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid authorization header format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "InvalidToken")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := mw.AuthMiddleware(func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("successful authentication", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := mw.AuthMiddleware(func(c echo.Context) error {
			adminID := c.Get("admin_id")
			assert.Equal(t, "admin-1", adminID)
			return c.String(http.StatusOK, "ok")
		})

		err := handler(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "ok", rec.Body.String())
	})
}

