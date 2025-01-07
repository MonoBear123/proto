// Пакет для инициализации и запуска приложения, включающего gRPC сервер и хранилище данных.
package app

import (
	grpcApp "auth/internal/app/grpc"
	"auth/internal/services/auth"
	"auth/internal/storage/postgres"
	"go.uber.org/zap"
	"time"
)

type App struct {
	GRPCServer *grpcApp.GRPC
	Storage    *postgres.Storage
}

// New - конструктор приложения, инициализирует компоненты приложения (gRPC сервер и хранилище).
//
// Параметры:
//   - log: Логгер для записи событий приложения.
//   - gRPCPort: Порт для gRPC сервера.
//   - port: Порт для подключения к базе данных.
//   - name: Имя базы данных.
//   - user: Имя пользователя для подключения к базе данных.
//   - password: Пароль для подключения к базе данных.
//   - tokenTTL: Время жизни токенов авторизации.
//
// Возвращает:
//   - Указатель на структуру App.
func New(log *zap.Logger,
	gRPCPort int,
	port int,
	name, user, password string,
	tokenTTL time.Duration,
) *App {
	storage := postgres.New(port, name, user, password)
	authServ := auth.New(log, storage, tokenTTL)
	GRPCApp := grpcApp.New(log, authServ, gRPCPort)

	return &App{
		GRPCServer: GRPCApp,
	}
}
