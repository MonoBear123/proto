package main

import (
	"client/internal/grpc/auth"
	"client/internal/grpc/predict"
	"client/internal/handlers/auth"
	"client/internal/handlers/predict"
	parser "client/internal/handlers/search"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	router := echo.New()
	router.Use(middleware.Logger())
	// Инициализация клиента gRPC-сервиса аутентификации.
	authGRPC := grpcAuth.New("auth_service:42022")
	// Инициализация клиента gRPC-сервиса предсказаний.
	predictGRPC := grpcPredict.New("predictor_service:42020")
	aHandler := authHandler.New(authGRPC)
	pHandler := predictHandler.New(predictGRPC)

	router.POST("/predict", pHandler.Predict)
	router.POST("/login", aHandler.Login, authHandler.ValidateDate)
	router.POST("/register", aHandler.Register, authHandler.ValidateDate)
	router.GET("/search", parser.Search)
	router.POST("/activate", aHandler.ActivateAccount)
	router.POST("/forgot-password", aHandler.ForgotPass)
	router.POST("/reset-password", aHandler.ResetPass)
	router.Logger.Fatal(router.Start(":8080"))
}
