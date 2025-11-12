package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("loads default values when env vars not set", func(t *testing.T) {
		// Clear all env vars
		os.Clearenv()
		
		cfg := Load()
		
		assert.Equal(t, 8080, cfg.Port)
		assert.Equal(t, "development", cfg.Env)
		assert.Equal(t, "localhost:50051", cfg.GRPCVenueAddr)
		assert.Equal(t, "localhost:50052", cfg.GRPCBookingAddr)
		assert.Equal(t, "localhost:6379", cfg.RedisAddr)
		assert.Equal(t, "", cfg.RedisPassword)
		assert.Equal(t, "change-me-in-production", cfg.JWTSecret)
		assert.Equal(t, "http://localhost:14268/api/traces", cfg.JaegerEndpoint)
	})
	
	t.Run("loads values from environment variables", func(t *testing.T) {
		os.Clearenv()
		os.Setenv("PORT", "9090")
		os.Setenv("ENV", "production")
		os.Setenv("GRPC_VENUE_ADDR", "venue:50051")
		os.Setenv("GRPC_BOOKING_ADDR", "booking:50052")
		os.Setenv("REDIS_ADDR", "redis:6379")
		os.Setenv("REDIS_PASSWORD", "secret123")
		os.Setenv("JWT_SECRET", "my-secret")
		os.Setenv("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces")
		
		cfg := Load()
		
		assert.Equal(t, 9090, cfg.Port)
		assert.Equal(t, "production", cfg.Env)
		assert.Equal(t, "venue:50051", cfg.GRPCVenueAddr)
		assert.Equal(t, "booking:50052", cfg.GRPCBookingAddr)
		assert.Equal(t, "redis:6379", cfg.RedisAddr)
		assert.Equal(t, "secret123", cfg.RedisPassword)
		assert.Equal(t, "my-secret", cfg.JWTSecret)
		assert.Equal(t, "http://jaeger:14268/api/traces", cfg.JaegerEndpoint)
		
		// Cleanup
		os.Clearenv()
	})
}

func TestGetEnv(t *testing.T) {
	t.Run("returns environment variable value when set", func(t *testing.T) {
		os.Setenv("TEST_VAR", "test-value")
		defer os.Unsetenv("TEST_VAR")
		
		result := getEnv("TEST_VAR", "default")
		
		assert.Equal(t, "test-value", result)
	})
	
	t.Run("returns default value when env var not set", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")
		
		result := getEnv("TEST_VAR", "default-value")
		
		assert.Equal(t, "default-value", result)
	})
	
	t.Run("returns default value when env var is empty", func(t *testing.T) {
		os.Setenv("TEST_VAR", "")
		defer os.Unsetenv("TEST_VAR")
		
		result := getEnv("TEST_VAR", "default-value")
		
		assert.Equal(t, "default-value", result)
	})
}

func TestGetEnvInt(t *testing.T) {
	t.Run("returns parsed integer from environment variable", func(t *testing.T) {
		os.Setenv("TEST_INT", "123")
		defer os.Unsetenv("TEST_INT")
		
		result := getEnvInt("TEST_INT", 0)
		
		assert.Equal(t, 123, result)
	})
	
	t.Run("returns default value when env var not set", func(t *testing.T) {
		os.Unsetenv("TEST_INT")
		
		result := getEnvInt("TEST_INT", 42)
		
		assert.Equal(t, 42, result)
	})
	
	t.Run("returns default value when env var is invalid", func(t *testing.T) {
		os.Setenv("TEST_INT", "not-a-number")
		defer os.Unsetenv("TEST_INT")
		
		result := getEnvInt("TEST_INT", 42)
		
		assert.Equal(t, 42, result)
	})
	
	t.Run("returns default value when env var is empty", func(t *testing.T) {
		os.Setenv("TEST_INT", "")
		defer os.Unsetenv("TEST_INT")
		
		result := getEnvInt("TEST_INT", 42)
		
		assert.Equal(t, 42, result)
	})
	
	t.Run("handles negative numbers", func(t *testing.T) {
		os.Setenv("TEST_INT", "-10")
		defer os.Unsetenv("TEST_INT")
		
		result := getEnvInt("TEST_INT", 0)
		
		assert.Equal(t, -10, result)
	})
}

func TestConfig_AllFields(t *testing.T) {
	t.Run("config struct has all required fields", func(t *testing.T) {
		cfg := &Config{
			Port:           8080,
			Env:            "test",
			GRPCVenueAddr:  "venue:50051",
			GRPCBookingAddr: "booking:50052",
			RedisAddr:      "redis:6379",
			RedisPassword:  "pass",
			JWTSecret:      "secret",
			JaegerEndpoint: "http://jaeger:14268/api/traces",
		}
		
		require.NotNil(t, cfg)
		assert.Equal(t, 8080, cfg.Port)
		assert.Equal(t, "test", cfg.Env)
		assert.Equal(t, "venue:50051", cfg.GRPCVenueAddr)
		assert.Equal(t, "booking:50052", cfg.GRPCBookingAddr)
		assert.Equal(t, "redis:6379", cfg.RedisAddr)
		assert.Equal(t, "pass", cfg.RedisPassword)
		assert.Equal(t, "secret", cfg.JWTSecret)
		assert.Equal(t, "http://jaeger:14268/api/traces", cfg.JaegerEndpoint)
	})
}

