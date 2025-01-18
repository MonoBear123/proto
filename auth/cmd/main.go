package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoadConfig()

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	application := app.New(logger, cfg.GRPC.Port, cfg.Storage.Port, cfg.Storage.Name,
		cfg.Storage.User, cfg.Storage.Password, cfg.TokenTTL)

	// Запуск сервера в отдельной горутине
	go application.GRPCServer.MustRun()

	logger.Info("application is started")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	application.GRPCServer.Stop()
	application.Storage.Close()

	logger.Info("application is stopped")

}
