package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

// main сборка и запуск сервиса авторизации
func main() {
	// Чтение конфига
	cfg := config.MustLoadConfig()

	// Создание объекта для логгирования
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	application := app.New(logger, cfg.GRPC.Port, cfg.Storage.Port, cfg.Storage.Name,
		cfg.Storage.User, cfg.Storage.Password, cfg.TokenTTL)

	// Запуск сервера в отдельной горутине
	go application.GRPCServer.MustRun()

	logger.Info("application is started")
	// Канал для отслеживания сигналов о завершения работы сервиса
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	// Ожидания сигнала в канале
	<-stop

	application.GRPCServer.Stop()
	application.Storage.Close()

	logger.Info("application is stopped")

}
