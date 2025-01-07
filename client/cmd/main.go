package main

import (
	"client/internal/grpc/auth"
	"client/internal/grpc/predict"
	"client/internal/handlers/auth"
	"client/internal/handlers/predict"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// main сборка и запуск клиент сервера.
func main() {

	router := echo.New()
	router.Use(middleware.Logger())
	// Инициализация клиента gRPC-сервиса аутентификации.
	authGRPC := grpcAuth.New("auth_service:42022")
	// Инициализация клиента gRPC-сервиса предсказаний.
	predictGRPC := grpcPredict.New("predictor_service:42020")
	// Создание обработчиков для маршрутов, связанных с аутентификацией.
	aHandler := authHandler.New(authGRPC)
	// Создание обработчиков для маршрутов, связанных с предсказаниями.
	pHandler := predictHandler.New(predictGRPC)

	// Определение маршрутов и их привязка к обработчикам.
	router.GET("/predict", pHandler.Predict)
	router.POST("/login", aHandler.Login)
	router.POST("/register", aHandler.Register)
	router.Logger.Fatal(router.Start(":8080"))
}
