package main

import (
	"client/internal/grpc/auth"
	grpcManager "client/internal/grpc/manage"
	"client/internal/grpc/predict"
	"client/internal/handlers/auth"
	"client/internal/handlers/manage"
	"client/internal/handlers/predict"
	parser "client/internal/handlers/search"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

func main() {

	router := echo.New()
	router.Use(middleware.Logger())
	authGRPC := grpcAuth.New("auth_service:42022")
	time.Sleep(5 * time.Second)
	predictGRPC := grpcPredict.New("predictor_service:42020")
	manageGRPC := grpcManager.New("auth_service:42022")
	mHandler := manageHandler.New(manageGRPC)
	aHandler := authHandler.New(authGRPC)
	pHandler := predictHandler.New(predictGRPC)
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	router.OPTIONS("/login", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	router.POST("/predict", pHandler.Predict)
	router.POST("/login", aHandler.Login, authHandler.ValidateDate)
	router.POST("/register", aHandler.Register, authHandler.ValidateDate)
	router.GET("/search", parser.Search)
	router.GET("/activate", mHandler.ActivateAccount)
	router.POST("/forgot-password", mHandler.ForgotPass)
	router.POST("/reset-password", mHandler.ResetPass)
	router.Logger.Fatal(router.Start(":8080"))
}
