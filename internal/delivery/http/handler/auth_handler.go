package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Auth handlers
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}

func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	ctx := c.Request().Context()
	if err := h.authUseCase.Register(ctx, req.Username, req.Password, req.Email); err != nil {
		if err.Error() == "username already exists" {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		if err.Error() == "username and password are required" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		log.Error().Err(err).Msg("Failed to register user")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message":  "User registered successfully",
		"username": req.Username,
	})
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	ctx := c.Request().Context()
	accessToken, refreshToken, err := h.authUseCase.Login(ctx, req.Username, req.Password)
	if err != nil {
		if err.Error() == "username and password are required" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		if err.Error() == "invalid credentials" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		log.Error().Err(err).Msg("Failed to login user")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()
	token, err := h.authUseCase.RefreshToken(ctx, "")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

