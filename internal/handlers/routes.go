package handlers

import (
	"database/sql"
	httpHandler "pvz_system/internal/handlers/http"
	middle "pvz_system/internal/middleware"

	"github.com/labstack/echo/v4"
)

func NewRouters(e *echo.Echo, db *sql.DB) {

	authHandler := httpHandler.NewUserHandler(db)
	pvzHandler := httpHandler.NewPVZHandler(db)
	receptionHandler := httpHandler.NewReceptionHandler(db)
	productHandler := httpHandler.NewProductHandler(db)

	jwtMiddleware := middle.JWTMiddleware(authHandler.AuthService)

	// Публичные маршруты
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)
	e.POST("/dummyLogin", authHandler.DummyLogin)

	// Группа с JWT аутентификацией
	authGroup := e.Group("")
	authGroup.Use(jwtMiddleware)

	// Маршруты только для модераторов
	adminGroup := authGroup.Group("")
	adminGroup.Use(middle.AdminOnlyMiddleware)
	adminGroup.POST("/pvz", pvzHandler.CreatePVZ)

	// Маршруты только для сотрудников
	employeeGroup := authGroup.Group("")
	employeeGroup.Use(middle.EmployeeOnlyMiddleware)
	employeeGroup.POST("/receptions", receptionHandler.StartReception)
	employeeGroup.POST("/products", productHandler.AddProduct)
	employeeGroup.POST("/pvz/:pvzId/close_last_reception", receptionHandler.CloseReception)
	employeeGroup.DELETE("/pvz/:pvzId/products/last", productHandler.RemoveLastProduct)

	// Общие маршруты для всех аутентифицированных
	authGroup.GET("/pvz", pvzHandler.ListPVZs)
	authGroup.GET("/pvz/:id", pvzHandler.GetPVZ)

}
