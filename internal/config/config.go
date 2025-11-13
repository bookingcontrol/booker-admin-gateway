package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port           int
	Env            string
	GRPCVenueAddr  string
	GRPCBookingAddr string
	RedisAddr      string
	RedisPassword  string
	JWTSecret      string
	OTLPEndpoint   string // OTLP endpoint (e.g., "localhost:4318")
}

func Load() *Config {
	return &Config{
		Port:           getEnvInt("PORT", 8080),
		Env:            getEnv("ENV", "development"),
		GRPCVenueAddr:  getEnv("GRPC_VENUE_ADDR", "localhost:50051"),
		GRPCBookingAddr: getEnv("GRPC_BOOKING_ADDR", "localhost:50052"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		OTLPEndpoint:   getEnv("OTLP_ENDPOINT", "localhost:4318"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

