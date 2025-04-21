package http

import (
	"net/http"
	"pvz_system/internal/models"
	"pvz_system/internal/services"

	"github.com/labstack/echo"
)

type UserHandler struct {
	authService services.AuthService
}

func NewUserHandler(authService services.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // client, moderator, employee
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DummyLoginRequest struct {
	Role string `json:"role"` // moderator, employee
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверный формат запроса"})
	}

	token, err := h.authService.Register(req.Email, req.Password, req.Role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, AuthResponse{Token: token})
}

func (h *UserHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверный формат запроса"})
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, AuthResponse{Token: token})
}

func (h *UserHandler) DummyLogin(c echo.Context) error {
	var req DummyLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверный формат запроса"})
	}

	role := models.UserRole(req.Role)
	token, err := h.authService.DummyLogin(c.Request().Context(), role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, AuthResponse{Token: token})
}
