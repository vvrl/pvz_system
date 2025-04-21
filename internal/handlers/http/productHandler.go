package http

import (
	"database/sql"
	"net/http"
	"pvz_system/internal/models"
	"pvz_system/internal/repository"
	"pvz_system/internal/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productService services.ProductService
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	repo := repository.NewProductRepository(db)
	receptonRepo := repository.NewReceptionRepository(db)
	service := services.NewProductService(repo, receptonRepo)
	return &ProductHandler{productService: service}
}

type AddProductRequest struct {
	Type models.ProductType `json:"type"` // Тип товара: electronic, clothes, shoes
}

type ProductResponse struct {
	ID          int                `json:"id"`
	ReceptionID int                `json:"reception_id"`
	Type        models.ProductType `json:"type"`
}

// POST /receptions/:id/products
func (h *ProductHandler) AddProduct(c echo.Context) error {
	receptionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Некорректный ID приемки"})
	}

	var req AddProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверный формат запроса"})
	}

	product, err := h.productService.AddProduct(c.Request().Context(), receptionID, req.Type)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, ProductResponse{
		ID:          product.ID,
		ReceptionID: product.ReceptionID,
		Type:        product.Type,
	})
}

// DELETE /receptions/:id/products/last
func (h *ProductHandler) RemoveLastProduct(c echo.Context) error {
	receptionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Некорректный ID приемки"})
	}

	if err := h.productService.RemoveLastProduct(c.Request().Context(), receptionID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GET /receptions/:id/products/last
func (h *ProductHandler) GetLastProduct(c echo.Context) error {
	receptionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Некорректный ID приемки"})
	}

	product, err := h.productService.GetLastProduct(c.Request().Context(), receptionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, ProductResponse{
		ID:          product.ID,
		ReceptionID: product.ReceptionID,
		Type:        product.Type,
	})
}
