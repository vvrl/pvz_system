package http

import (
	"database/sql"
	"net/http"
	"pvz_system/internal/repository"
	"pvz_system/internal/services"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type PVZHandler struct {
	pvzService services.PVZService
}

func NewPVZHandler(db *sql.DB) *PVZHandler {
	repo := repository.NewPVZRepository(db)
	svc := services.NewPVZService(repo)
	return &PVZHandler{pvzService: svc}
}

func (h *PVZHandler) CreatePVZ(c echo.Context) error {
	type request struct {
		City string `json:"city"`
	}

	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	pvz, err := h.pvzService.CreatePVZ(c.Request().Context(), req.City)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, pvz)
}

func (h *PVZHandler) GetPVZ(c echo.Context) error {
	id := c.Param("id")

	pvzID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid PVZ ID"})
	}

	pvz, err := h.pvzService.GetPVZ(c.Request().Context(), pvzID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, pvz)
}

func (h *PVZHandler) ListPVZs(c echo.Context) error {
	startDateStr := c.QueryParam("startDate")
	endDateStr := c.QueryParam("endDate")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	startDate, _ := time.Parse(time.RFC3339, startDateStr)
	endDate, _ := time.Parse(time.RFC3339, endDateStr)
	page := parseIntOrDefault(pageStr, 1)
	limit := parseIntOrDefault(limitStr, 10)

	pvzs, err := h.pvzService.ListPVZs(c.Request().Context(), startDate, endDate, page, limit)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, pvzs)
}

func parseIntOrDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return def
}
