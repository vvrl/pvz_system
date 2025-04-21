package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pvz_system/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := "INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return user, err
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, email, password, role FROM users WHERE email = $1"

	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("ошибка вывода пользователя: %w", err)
	}

	return &user, nil
}
