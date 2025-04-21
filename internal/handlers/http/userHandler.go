package http

import (
	"database/sql"
	"net/http"
	"pvz_system/internal/models"
	"pvz_system/internal/repository"
	"pvz_system/internal/services"
	"time"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	AuthService services.AuthService
}

func NewUserHandler(db *sql.DB) *UserHandler {
	repo := repository.NewUserRepository(db)
	service := services.NewAuthService("KEY", time.Hour*24, repo)

	return &UserHandler{AuthService: *service}
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

	token, err := h.AuthService.Register(req.Email, req.Password, req.Role)
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

	token, err := h.AuthService.Login(req.Email, req.Password)
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
	token, err := h.AuthService.DummyLogin(c.Request().Context(), role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, AuthResponse{Token: token})
}
