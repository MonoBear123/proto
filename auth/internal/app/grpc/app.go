// Пакет для инициализации и звпуска  gRPC сервера.
// Включает в себя создание сервера, запуск, а также его завершение.
package grpcApp

import (
	authgrpc "auth/internal/grpc/auth"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GRPC struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

// New - конструктор для создания нового gRPC сервера.
//
// Параметры:
//   - log: Логгер для записи информации о сервере.
//   - authServ: Экземпляр службы аутентификации для регистрации на сервере.
//   - port: Порт для прослушивания gRPC сервера.
//
// Возвращает:
//   - Указатель на структуру GRPC с инициализированным сервером.
func New(log *zap.Logger, authServ authgrpc.AuthService, port int) *GRPC {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer, authServ)
	return &GRPC{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun - вызывает функцию Run.
// При возникновении ошибки вызывается panic.
func (g *GRPC) MustRun() {
	if err := g.Run(); err != nil {
		panic(err)
	}
}

// Run - запускает gRPC сервер.
// Использует для запуска порт, который был передан при создании структуры GRPC.
func (g *GRPC) Run() error {
	const op = "grpcApp.Run"
	log := g.log.With(zap.String("op", op), zap.Int("port", g.port))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("starting gRPC server")

	if err := g.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop - останавливает gRPC сервер с использованием graceful shutdown.
func (g *GRPC) Stop() {
	const op = "grpcApp.Stop"
	log := g.log.With(zap.String("op", op))
	log.Info("stopping gRPC server")

	g.gRPCServer.GracefulStop() //  GracefulStop позволяет корректно завершить все активные соединения.
}
