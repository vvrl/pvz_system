package services

import (
	"context"
	"errors"
	"fmt"
	"pvz_system/internal/models"
	"pvz_system/internal/repository"
)

type productService struct {
	productRepo   repository.ProductRepository
	receptionRepo repository.ReceptionRepository
}

func NewProductService(productRepo repository.ProductRepository, receptionRepo repository.ReceptionRepository) *productService {
	return &productService{productRepo: productRepo,
		receptionRepo: receptionRepo,
	}
}

type ProductService interface {
	AddProduct(ctx context.Context, receptionID int, productType models.ProductType) (*models.Product, error)
	RemoveLastProduct(ctx context.Context, receptionID int) error
	GetLastProduct(ctx context.Context, receptionID int) (*models.Product, error)
}

func (s *productService) AddProduct(ctx context.Context, receptionID int, productType models.ProductType) (*models.Product, error) {
	// проверяем тип товара
	if !isValidProductType(productType) {
		return nil, errors.New("некорректный тип товара")
	}

	// проверяем статус приемки
	reception, err := s.receptionRepo.GetReceptionByID(ctx, receptionID)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки приемки: %w", err)
	}

	if reception.Status == models.StatusClosed {
		return nil, errors.New("нельзя добавить товар в закрытую приемку")
	}

	// создаем товар
	product := &models.Product{
		ReceptionID: receptionID,
		Type:        productType,
	}

	createdProduct, err := s.productRepo.AddProduct(ctx, receptionID, product)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания товара: %w", err)
	}

	return createdProduct, nil
}

func (s *productService) GetLastProduct(ctx context.Context, receptionID int) (*models.Product, error) {
	product, err := s.productRepo.GetLastProduct(ctx, receptionID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения товара: %w", err)
	}
	return product, nil
}

func (s *productService) RemoveLastProduct(ctx context.Context, receptionID int) error {
	// проверяем статус приемки
	reception, err := s.receptionRepo.GetReceptionByID(ctx, receptionID)
	if err != nil {
		return fmt.Errorf("ошибка проверки приемки: %w", err)
	}

	if reception.Status == models.StatusClosed {
		return errors.New("примка уже закрыта")
	}

	// получаем последний товар
	lastProduct, err := s.productRepo.GetLastProduct(ctx, receptionID)
	if err != nil {
		return fmt.Errorf("ошибка получения последнего товара: %w", err)
	}

	// удаляем товар
	if err := s.productRepo.RemoveProduct(ctx, lastProduct.ID); err != nil {
		return fmt.Errorf("ошибка удаления товара: %w", err)
	}

	return nil
}

func isValidProductType(productType models.ProductType) bool {
	switch productType {
	case models.TypeElectronic, models.TypeClothes, models.TypeShoes:
		return true
	default:
		return false
	}
}
