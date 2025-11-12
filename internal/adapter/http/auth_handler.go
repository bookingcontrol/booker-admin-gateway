package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	uc "github.com/bookingcontrol/booker-admin-gateway/internal/usecase/auth"
)

type AuthHandler struct {
	svc *uc.Service
}

func NewAuthHandler(svc *uc.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type registerReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email,omitempty"`
}

type loginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req registerReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	out, err := h.svc.Register(c.Request().Context(), uc.CreateInput{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		if err.Error() == "username already exists" {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		if err.Error() == "username and password are required" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		log.Error().Err(err).Msg("Failed to register user")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, out)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req loginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	out, err := h.svc.Login(c.Request().Context(), uc.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
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

	return c.JSON(http.StatusOK, out)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	token, err := h.svc.RefreshToken(c.Request().Context(), "")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

