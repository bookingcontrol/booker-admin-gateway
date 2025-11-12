package redis

import (
	"context"

	dom "github.com/bookingcontrol/booker-admin-gateway/internal/domain/auth"
	"github.com/bookingcontrol/booker-admin-gateway/internal/infrastructure/redis"
)

type AuthRepo struct {
	client *redis.Client
}

func NewAuthRepo(client *redis.Client) dom.Repository {
	return &AuthRepo{
		client: client,
	}
}

func (r *AuthRepo) UserExists(ctx context.Context, username string) (bool, error) {
	exists, err := r.client.Exists(ctx, "user:"+username)
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *AuthRepo) GetUserPassword(ctx context.Context, username string) (string, error) {
	return r.client.HGet(ctx, "user:"+username, "password")
}

func (r *AuthRepo) CreateUser(ctx context.Context, username string, userData map[string]interface{}) error {
	return r.client.HSet(ctx, "user:"+username, userData)
}

