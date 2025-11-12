package repository

import "context"

// UserRepository defines interface for user storage operations
type UserRepository interface {
	UserExists(ctx context.Context, username string) (bool, error)
	GetUserPassword(ctx context.Context, username string) (string, error)
	CreateUser(ctx context.Context, username string, userData map[string]interface{}) error
}
