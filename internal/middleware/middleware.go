package middleware

import (
	"errors"
	"net/http"
	"strings"

	"pvz_system/internal/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	authHeader   = "Authorization"
	bearerPrefix = "Bearer "
)

// JWTClaims - кастомные claims для нашего JWT токена
type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type ErrorResponse struct {
	Message string      `json:"message"`           // Сообщение об ошибке
	Code    string      `json:"code,omitempty"`    // Опциональный код ошибки (для классификации)
	Details interface{} `json:"details,omitempty"` // Дополнительные детали ошибки
}

func JWTMiddleware(authService services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// извлекаем токен из заголовка
			tokenString, err := extractToken(c.Request())
			if err != nil {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
			}

			// валидируем токен
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid token"})
			}

			// добавляем claims в контекст
			c.Set("user_id", claims.UserID)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get(authHeader)
	if authHeader == "" {
		return "", errors.New("требуется заголовок авторизации")
	}

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("формат заголовка авторизации должен быть 'Bearer {token}'")
	}

	return authHeader[len(bearerPrefix):], nil
}

func AdminOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("user_role").(string)
		if role != "moderator" {
			return c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "Доступ запрещен. Требуются права модератора",
				Code:    "forbidden",
			})
		}
		return next(c)
	}
}

func EmployeeOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := c.Get("user_role").(string)
		if role != "employee" {
			return c.JSON(http.StatusForbidden, ErrorResponse{
				Message: "Доступ запрещен. Необходимые права сотрудников",
				Code:    "forbidden",
			})
		}
		return next(c)
	}
}
