package handlers

import (
	"database/sql"
	httpHandler "pvz_system/internal/handlers/http"
	middle "pvz_system/internal/middleware"

	"github.com/labstack/echo"
)

func NewRouters(e *echo.Echo, db *sql.DB) {
	//jwtMiddleware := middle.JWTMiddleware(authService)

	authHandler := httpHandler.NewUserHandler(authService)
	pvzHandler := httpHandler.NewPVZHandler(db)
	receptionHandler := httpHandler.NewReceptionHandler(db)

	// Публичные маршруты
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)
	e.POST("/dummyLogin", authHandler.DummyLogin)

	// Группа с JWT аутентификацией
	authGroup := e.Group("")
	authGroup.Use(middle.jwtMiddleware)

	// Маршруты только для модераторов
	adminGroup := authGroup.Group("")
	adminGroup.Use(middle.AdminOnlyMiddleware)
	adminGroup.POST("/pvz", httpHandler.PVZHandler.CreatePVZ)

	// Маршруты только для сотрудников
	employeeGroup := authGroup.Group("")
	employeeGroup.Use(EmployeeOnlyMiddleware)
	employeeGroup.POST("/receptions", receptionHandler.CreateReception)
	employeeGroup.POST("/products", productHandler.AddProduct)
	employeeGroup.POST("/pvz/:pvzId/close_last_reception", pvzHandler.CloseLastReception)
	employeeGroup.DELETE("/pvz/:pvzId/products/last", productHandler.RemoveLastProduct)

	// Общие маршруты для всех аутентифицированных
	authGroup.GET("/pvz", pvzHandler.ListPVZs)
	authGroup.GET("/pvz/:id", pvzHandler.GetPVZ)

}
