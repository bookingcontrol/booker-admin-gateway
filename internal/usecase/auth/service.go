package auth

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	dom "github.com/bookingcontrol/booker-admin-gateway/internal/domain/auth"
)

type Service struct {
	repo dom.Repository
}

func NewService(repo dom.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, in CreateInput) (RegisterView, error) {
	if in.Username == "" || in.Password == "" {
		return RegisterView{}, errors.New("username and password are required")
	}

	exists, err := s.repo.UserExists(ctx, in.Username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check user existence")
		return RegisterView{}, errors.New("internal server error")
	}
	if exists {
		return RegisterView{}, errors.New("username already exists")
	}

	userData := map[string]interface{}{
		"username": in.Username,
		"password": in.Password, // TODO: Hash password
		"email":    in.Email,
	}

	if err := s.repo.CreateUser(ctx, in.Username, userData); err != nil {
		log.Error().Err(err).Msg("Failed to store user")
		return RegisterView{}, errors.New("failed to create user")
	}

	log.Info().Str("username", in.Username).Msg("User registered")
	return RegisterView{
		Username: in.Username,
		Message:  "User registered successfully",
	}, nil
}

func (s *Service) Login(ctx context.Context, in LoginInput) (LoginView, error) {
	if in.Username == "" || in.Password == "" {
		return LoginView{}, errors.New("username and password are required")
	}

	exists, err := s.repo.UserExists(ctx, in.Username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check user existence")
		return LoginView{}, errors.New("internal server error")
	}
	if !exists {
		return LoginView{}, errors.New("invalid credentials")
	}

	storedPassword, err := s.repo.GetUserPassword(ctx, in.Username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user password")
		return LoginView{}, errors.New("internal server error")
	}

	if storedPassword != in.Password {
		return LoginView{}, errors.New("invalid credentials")
	}

	// Generate token (in production, use JWT)
	token := "token-" + in.Username         // TODO: Generate proper JWT
	refreshToken := "refresh-" + in.Username // TODO: Generate proper refresh token

	log.Info().Str("username", in.Username).Msg("User logged in")
	return LoginView{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// TODO: Implement token refresh
	return "dummy-token", nil
}

