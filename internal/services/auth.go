package services

import (
	"context"
	"errors"
	"fmt"
	"pvz_system/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	secretKey     string
	tokenDuration time.Duration
	userRepo      repository.UserRepository
}

type Claims struct {
	UserID int    `json:"userId"`
	Role   string `json:"role"` // Пример: можно добавить роль
	jwt.RegisteredClaims
}

func NewAuthService(secretKey string, tokenDuration time.Duration, userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
		userRepo:      userRepo,
	}
}

func (s *AuthService) generateToken(userID int, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pvz_system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	// Парсинг токена с проверкой подписи
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("некоректный токен")
		}
		return []byte(s.secretKey), nil
	},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("просроченный токен: %w", err) // Отдельная ошибка для просроченного токена
		}
		return nil, fmt.Errorf("некоректный токен: %w", err)
	}

	// Проверка валидности токена и типа claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if _, err := s.userRepo.GetUserById(context.Background(), claims.UserID); err != nil {
			return nil, fmt.Errorf("пользователь не найден: %w", err)
		}
		return claims, nil // Успешная валидация
	}

	return nil, fmt.Errorf("некоректный токен: %w", err)
}
