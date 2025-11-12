package repository

import (
	"context"

	"github.com/bookingcontrol/booker-admin-gateway/internal/domain/repository"
	"github.com/bookingcontrol/booker-admin-gateway/internal/infrastructure/redis"
)

type userRepository struct {
	redisClient *redis.Client
}

func NewUserRepository(redisClient *redis.Client) repository.UserRepository {
	return &userRepository{
		redisClient: redisClient,
	}
}

func (r *userRepository) UserExists(ctx context.Context, username string) (bool, error) {
	exists, err := r.redisClient.Exists(ctx, "user:"+username)
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *userRepository) GetUserPassword(ctx context.Context, username string) (string, error) {
	return r.redisClient.HGet(ctx, "user:"+username, "password")
}

func (r *userRepository) CreateUser(ctx context.Context, username string, userData map[string]interface{}) error {
	return r.redisClient.HSet(ctx, "user:"+username, userData)
}
