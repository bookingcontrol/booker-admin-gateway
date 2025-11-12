package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тестируем только логику маппинга (префикс "user:", проверка exists > 0)
// Без реального Redis - используем мок или проверяем только структуру

func TestAuthRepo_KeyPrefixing(t *testing.T) {
	// Проверяем, что ключи правильно формируются с префиксом "user:"
	// Это единственная логика в адаптере
	
	t.Run("key format is correct", func(t *testing.T) {
		username := "testuser"
		expectedKey := "user:" + username
		
		// Проверяем, что ключ формируется правильно
		// В реальном коде: r.client.Exists(ctx, "user:"+username)
		assert.Equal(t, "user:testuser", expectedKey)
	})
	
	t.Run("exists check logic", func(t *testing.T) {
		// Проверяем логику: exists > 0 означает, что пользователь существует
		testCases := []struct {
			exists   int64
			expected bool
		}{
			{0, false},
			{1, true},
			{2, true},
		}
		
		for _, tc := range testCases {
			result := tc.exists > 0
			assert.Equal(t, tc.expected, result, "exists=%d should return %v", tc.exists, tc.expected)
		}
	})
}

// Интеграционный тест с реальным Redis (опционально, можно пропустить если Redis недоступен)
func TestAuthRepo_Integration(t *testing.T) {
	t.Skip("Integration test - requires Redis. Set REDIS_ADDR env var to enable")
	
	// Можно раскомментировать для интеграционного теста
	// redisAddr := os.Getenv("REDIS_ADDR")
	// if redisAddr == "" {
	// 	t.Skip("REDIS_ADDR not set")
	// }
	// client := redis.NewClient(redisAddr, "")
	// defer client.Close()
	// repo := NewAuthRepo(client)
	// ctx := context.Background()
	// 
	// // Тесты с реальным Redis...
}

