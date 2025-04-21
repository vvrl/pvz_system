package services

import (
	"context"
	"errors"
	"fmt"
	"pvz_system/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Register(email, password, role string) (string, error) {
	// Валидация роли
	if role != "client" && role != "moderator" && role != "employee" {
		return "", fmt.Errorf("invalid role")
	}

	// Хеширование пароля
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Создание пользователя
	user, err := s.userRepo.CreateUser(context.Background(), email, hashedPassword, models.UserRole(role))
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	// Генерация токена
	return s.generateToken(user.ID, string(user.Role))
}

func (s *AuthService) Login(email, password string) (string, error) {
	// Получение пользователя
	user, err := s.userRepo.GetUserByEmail(context.Background(), email)
	if err != nil {
		return "", errors.New("некорректный email или пароль")
	}

	// Проверка пароля
	if !checkPasswordHash(password, user.Password) {
		return "", errors.New("некорректный email или пароль")
	}

	// Генерация токена
	return s.generateToken(user.ID, string(user.Role))
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
