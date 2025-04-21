package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pvz_system/internal/models"
	"time"
)

type ReceptionRepository interface {
	CreateReception(ctx context.Context, pvzID int) (*models.Reception, error)
	GetReceptionByID(ctx context.Context, id int) (*models.Reception, error)
	GetOpenReception(ctx context.Context, pvzID int) (*models.Reception, error)
	CloseReception(ctx context.Context, id int) error
}

type ReceptionRepo struct {
	db *sql.DB
}

func NewReceptionRepository(db *sql.DB) ReceptionRepository {
	return &ReceptionRepo{db}
}

func (r *ReceptionRepo) CreateReception(ctx context.Context, pvzID int) (*models.Reception, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия транзакции: %w", err)
	}
	defer tx.Rollback()

	// Проверяем, есть ли открытая приёмка
	var openReceiptID int
	query := "SELECT id FROM receptions WHERE pvzId = $1 AND status = $2"
	err = tx.QueryRowContext(ctx, query, pvzID, models.StatusInProgress).Scan(&openReceiptID)

	if err == nil {
		return nil, fmt.Errorf("уже есть открытая приемка: %w", err)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("ошибка проверки открытых приемок: %w", err)
	}
	var receptionID int
	query = "INSERT INTO receptions (pvzId, dateTime, status) VALUES ($1, $2, $3) RETURNING id"
	err = tx.QueryRowContext(ctx, query, pvzID, time.Now(), models.StatusInProgress).Scan(&receptionID)

	if err != nil {
		return nil, fmt.Errorf("неудачнач поптыка создания приемки: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("ошибка при закрытии транзакции: %w", err)
	}

	return &models.Reception{
		ID:       receptionID,
		PVZID:    pvzID,
		DateTime: time.Now(),
		Status:   models.StatusInProgress,
	}, nil

}

func (r *ReceptionRepo) GetReceptionByID(ctx context.Context, id int) (*models.Reception, error) {
	query := "SELECT id, dateTime, pvzId, status FROM receptions WHERE id = $1"

	var reception models.Reception

	err := r.db.QueryRowContext(ctx, query, id).Scan(&reception.ID, &reception.PVZID, &reception.DateTime, &reception.Status)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения приемки: %w", err)
	}

	return &reception, nil
}

func (r *ReceptionRepo) GetOpenReception(ctx context.Context, pvzID int) (*models.Reception, error) {
	query := "SELECT id, dateTime, pvzId, status FROM receptions WHERE pvzId = $1 AND status = $2 ORDER BY id DESC LIMIT 1"
	var reception models.Reception

	err := r.db.QueryRowContext(ctx, query, pvzID, models.StatusInProgress).Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status)
	if err != nil {
		return nil, fmt.Errorf("ощибка в получении последней незакрытой приемки: %w", err)
	}

	return &reception, nil
}

func (r *ReceptionRepo) CloseReception(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	// Проверяем статус приёмки
	var status string
	query := "SELECT status FROM receptions WHERE id = $1 FOR UPDATE"
	err = tx.QueryRowContext(ctx, query, id).Scan(&status)

	if err != nil {
		return fmt.Errorf("ошибка получения статуса приемки: %w", err)
	}

	if status == string(models.StatusClosed) {
		return fmt.Errorf("приемка уже закрыта: %w", err)
	}

	// Закрываем приёмку

	query = "UPDATE receptions SET status = $1 WHERE id = $2"
	_, err = tx.ExecContext(ctx, string(models.StatusClosed), id)
	if err != nil {
		return fmt.Errorf("failed to close receipt: %w", err)
	}

	return tx.Commit()
}
