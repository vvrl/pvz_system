package http

import (
	"database/sql"
	"errors"
	"net/http"
	"pvz_system/internal/repository"
	"pvz_system/internal/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ReceptionHandler struct {
	receptionService services.ReceptionService
}

func NewReceptionHandler(db *sql.DB) *ReceptionHandler {
	repo := repository.NewReceptionRepository(db)
	pvzRepo := repository.NewPVZRepository(db)
	svc := services.NewReceptiontService(repo, pvzRepo)
	return &ReceptionHandler{receptionService: svc}
}

type CreateReceptionRequest struct {
	PVZID int `json:"pvzId"`
}

func (h *ReceptionHandler) StartReception(c echo.Context) error {
	var req CreateReceptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверный формат запроса"})
	}

	reception, err := h.receptionService.StartReception(c.Request().Context(), req.PVZID)
	if err != nil {
		if err.Error() == "уже есть открытая приемка" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, reception)
}

func (h *ReceptionHandler) CloseReception(c echo.Context) error {
	pvzIDStr := c.Param("pvzId")
	pvzID, err := strconv.Atoi(pvzIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверный ID ПВЗ"})
	}

	openReception, err := h.receptionService.GetOpenReception(c.Request().Context(), pvzID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Открытая приемка не найдена"})
	}

	err = h.receptionService.CloseReception(c.Request().Context(), openReception.ID)
	if err != nil {
		if errors.Is(err, errors.New("приемка уже закрыта")) {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, openReception)
}
