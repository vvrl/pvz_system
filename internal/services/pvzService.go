package services

import (
	"context"
	"errors"
	"fmt"
	"pvz_system/internal/models"
	"pvz_system/internal/repository"
	"time"
)

type pvzService struct {
	pvzRepo repository.PvzRepository
}

func NewPVZService(pvzRepo repository.PvzRepository) *pvzService {
	return &pvzService{pvzRepo: pvzRepo}
}

type PVZService interface {
	CreatePVZ(ctx context.Context, city string) (*models.PVZ, error)
	GetPVZ(ctx context.Context, id int) (*models.PVZ, error)
	ListPVZs(ctx context.Context, startDate, endDate time.Time, page, pageSize int) ([]models.PVZ, error)
}

func (s *pvzService) CreatePVZ(ctx context.Context, city string) (*models.PVZ, error) {
	allowedCities := map[string]bool{
		"Москва":          true,
		"Санкт-Петербург": true,
		"Казань":          true,
	}

	if !allowedCities[city] {
		return nil, errors.New("недоступный город для ПВЗ")
	}

	pvz := &models.PVZ{
		City:             city,
		RegistrationDate: time.Now(),
	}

	createdPVZ, err := s.pvzRepo.CreatePVZ(ctx, pvz)
	if err != nil {
		return nil, fmt.Errorf("ошибка присоздании ПВЗ: %w", err)
	}

	return createdPVZ, nil
}

func (s *pvzService) GetPVZ(ctx context.Context, id int) (*models.PVZ, error) {
	pvz, err := s.pvzRepo.GetPVZByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get PVZ: %w", err)
	}
	return pvz, nil
}

func (s *pvzService) ListPVZs(ctx context.Context, startDate, endDate time.Time, page, pageSize int) ([]models.PVZ, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 30 {
		pageSize = 10 // Устанавливаем разумный лимит по умолчанию
	}

	// Проверка временного диапазона
	if startDate.After(endDate) {
		return nil, errors.New("некорректный временной диапазон")
	}

	// Рассчитываем offset для пагинации
	offset := (page - 1) * pageSize

	// Вызываем репозиторий
	pvzs, err := s.pvzRepo.ListPVZs(ctx, startDate, endDate, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка ПВЗ: %w", err)
	}

	return pvzs, nil
}
