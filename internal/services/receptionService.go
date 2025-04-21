package services

import (
	"context"
	"errors"
	"fmt"
	"pvz_system/internal/models"
	"pvz_system/internal/repository"
	"strings"
)

type receptionService struct {
	receptionRepo repository.ReceptionRepository
	pvzRepo       repository.PvzRepository
}

func NewReceptiontService(receiptRepo repository.ReceptionRepository, pvzRepo repository.PvzRepository) *receptionService {
	return &receptionService{receptionRepo: receiptRepo, pvzRepo: pvzRepo}
}

type ReceptionService interface {
	StartReception(ctx context.Context, pvzID int) (*models.Reception, error)
	CloseReception(ctx context.Context, id int) error
	GetReception(ctx context.Context, id int) (*models.Reception, error)
	GetOpenReception(ctx context.Context, pvzID int) (*models.Reception, error)
}

func (s *receptionService) StartReception(ctx context.Context, pvzID int) (*models.Reception, error) {
	if _, err := s.pvzRepo.GetPVZByID(ctx, pvzID); err != nil {
		return nil, fmt.Errorf("%w: %v", errors.New("ПВЗ не найден"), err)
	}

	reception, err := s.receptionRepo.CreateReception(ctx, pvzID)
	if err != nil {
		if strings.HasPrefix(string(err.Error()), "уже есть открытая приемка") {
			return nil, errors.New("уже есть открытая приемка")
		}
		return nil, fmt.Errorf("ошибка создания приемки: %w", err)
	}

	return reception, nil
}

func (s *receptionService) CloseReception(ctx context.Context, id int) error {
	if _, err := s.receptionRepo.GetReceptionByID(ctx, id); err != nil {

		return fmt.Errorf("ошибка проверки приемки: %w", err)
	}

	if err := s.receptionRepo.CloseReception(ctx, id); err != nil {
		if strings.HasPrefix(string(err.Error()), "приемка уже закрыта") {
			return errors.New("приемка уже закрыта")
		}
		return fmt.Errorf("ошибка закрытия приемки: %w", err)
	}

	return nil
}

func (s *receptionService) GetReception(ctx context.Context, id int) (*models.Reception, error) {
	reception, err := s.receptionRepo.GetReceptionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения приемки: %w", err)
	}
	return reception, nil
}
func (s *receptionService) GetOpenReception(ctx context.Context, pvzID int) (*models.Reception, error) {
	reception, err := s.receptionRepo.GetOpenReception(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения открытой приемки: %w", err)
	}
	return reception, nil
}
