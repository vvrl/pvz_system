package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pvz_system/internal/models"
	"time"
)

type ProductRepository interface {
	AddProduct(ctx context.Context, receptionID int, product *models.Product) (*models.Product, error)
	GetLastProduct(ctx context.Context, receptionID int) (*models.Product, error)
	RemoveProduct(ctx context.Context, productID int) error
}

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &ProductRepo{db}
}

func (r *ProductRepo) AddProduct(ctx context.Context, receptionID int, product *models.Product) (*models.Product, error) {
	query := "INSERT INTO products (receptionId, type, dateTime) VALUES ($1, $2, $3) RETURNING id"
	var id int

	err := r.db.QueryRowContext(ctx, query, receptionID, product.Type, time.Now()).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания предмета: %w", err)
	}

	product.ID = id
	product.DateTime = time.Now()

	return product, nil
}

func (r *ProductRepo) GetLastProduct(ctx context.Context, receptionID int) (*models.Product, error) {
	query := `
		SELECT id, receptionId, type, dateTime FROM products 
		WHERE receptionId = $1
		ORDER BY dateTime DESC
		LIMIT 1
	`

	var product models.Product
	err := r.db.QueryRowContext(ctx, query, receptionID).Scan(
		&product.ID,
		&product.ReceptionID,
		&product.Type,
		&product.DateTime,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении последнего предмета: %w", err)
	}

	return &product, nil
}

func (r *ProductRepo) RemoveProduct(ctx context.Context, productID int) error {

	query := "DELETE FROM products WHERE id = $1`"

	result, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении предмета: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка приполучении затронутых строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("запрос на удаление предмета не сработал: %w", err)
	}

	return nil
}
