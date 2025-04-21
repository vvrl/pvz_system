package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pvz_system/internal/models"
	"time"
)

type PvzRepository interface {
	CreatePVZ(ctx context.Context, pvz *models.PVZ) (*models.PVZ, error)
	GetPVZByID(ctx context.Context, id int) (*models.PVZ, error)
	ListPVZs(ctx context.Context, fromTime, toTime time.Time, limit, offset int) ([]models.PVZ, error)
}

type PVZRepo struct {
	db *sql.DB
}

func NewPVZRepository(db *sql.DB) PvzRepository {
	return &PVZRepo{db}
}

func (r *PVZRepo) CreatePVZ(ctx context.Context, pvz *models.PVZ) (*models.PVZ, error) {
	query := "INSERT INTO pvz (city, registrationDate) VALUES ($1, $2) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, pvz.City, pvz.RegistrationDate).Scan(&pvz.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания ПВЗ: %w", err)
	}
	return pvz, nil
}

func (r *PVZRepo) GetPVZByID(ctx context.Context, id int) (*models.PVZ, error) {
	query := "SELECT id, city, registrationDate FROM pvz WHERE id = $1"

	var pvz models.PVZ
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&pvz.ID,
		&pvz.City,
		&pvz.RegistrationDate,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ПВЗ: %w", err)
	}

	return &pvz, nil
}

func (r *PVZRepo) ListPVZs(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]models.PVZ, error) {
	query := `
        SELECT pvz.id, pvz.city, pvz.registrationDate,
            r.id, r.dateTime, r.pvzId, r.status,
        FROM pvz
        JOIN receptions r ON pvz.id = r.pvzId
        WHERE r.dateTime BETWEEN $1 AND $2
        ORDER BY pvz.id, r.dateTime
		LIMIT $3 OFFSET $4
    `
	rows, err := r.db.QueryContext(ctx, query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе списка пвз: %w", err)
	}

	defer rows.Close()
	var pvzs []models.PVZ // Итоговый список ПВЗ

	for rows.Next() {
		var currentPVZ models.PVZ
		if err := rows.Scan(&currentPVZ.ID, &currentPVZ.City, &currentPVZ.RegistrationDate); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строк: %w", err)
		}

		pvzs = append(pvzs, currentPVZ)
	}
	return pvzs, nil
}
