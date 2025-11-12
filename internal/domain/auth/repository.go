package auth

import "context"

// Repository defines interface for user storage operations
type Repository interface {
	UserExists(ctx context.Context, username string) (bool, error)
	GetUserPassword(ctx context.Context, username string) (string, error)
	CreateUser(ctx context.Context, username string, userData map[string]interface{}) error
}

