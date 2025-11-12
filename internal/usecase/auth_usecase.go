package usecase

import (
	"context"
	"errors"

	"github.com/bookingcontrol/booker-admin-gateway/internal/domain/repository"
	"github.com/rs/zerolog/log"
)

type authUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(userRepo repository.UserRepository) *authUseCase {
	return &authUseCase{
		userRepo: userRepo,
	}
}

func (uc *authUseCase) Register(ctx context.Context, username, password, email string) error {
	if username == "" || password == "" {
		return errors.New("username and password are required")
	}

	exists, err := uc.userRepo.UserExists(ctx, username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check user existence")
		return errors.New("internal server error")
	}
	if exists {
		return errors.New("username already exists")
	}

	userData := map[string]interface{}{
		"username": username,
		"password": password, // TODO: Hash password
		"email":    email,
	}

	if err := uc.userRepo.CreateUser(ctx, username, userData); err != nil {
		log.Error().Err(err).Msg("Failed to store user")
		return errors.New("failed to create user")
	}

	log.Info().Str("username", username).Msg("User registered")
	return nil
}

func (uc *authUseCase) Login(ctx context.Context, username, password string) (accessToken, refreshToken string, err error) {
	if username == "" || password == "" {
		return "", "", errors.New("username and password are required")
	}

	exists, err := uc.userRepo.UserExists(ctx, username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check user existence")
		return "", "", errors.New("internal server error")
	}
	if !exists {
		return "", "", errors.New("invalid credentials")
	}

	storedPassword, err := uc.userRepo.GetUserPassword(ctx, username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user password")
		return "", "", errors.New("internal server error")
	}

	if storedPassword != password {
		return "", "", errors.New("invalid credentials")
	}

	// Generate token (in production, use JWT)
	token := "token-" + username         // TODO: Generate proper JWT
	refreshToken = "refresh-" + username // TODO: Generate proper refresh token

	log.Info().Str("username", username).Msg("User logged in")
	return token, refreshToken, nil
}

func (uc *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// TODO: Implement token refresh
	return "dummy-token", nil
}
