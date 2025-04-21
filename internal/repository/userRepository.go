package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pvz_system/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, password string, role models.UserRole) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) CreateUser(ctx context.Context, email, password string, role models.UserRole) (*models.User, error) {
	query := "INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id"
	var id int
	err := r.db.QueryRowContext(ctx, query, email, password, role).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return &models.User{
		ID:    id,
		Email: email,
		Role:  role,
	}, err
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, email, password, role FROM users WHERE email = $1"

	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("ошибка вывода пользователя по почте: %w", err)
	}

	return &user, nil
}

func (r *userRepo) GetUserById(ctx context.Context, id int) (*models.User, error) {
	query := "SELECT id, email, password, role FROM users WHERE id = $1"

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("ошибка вывода пользователя по id: %w", err)
	}

	return &user, nil
}
