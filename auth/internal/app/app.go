package app

import (
	grpcApp "auth/internal/app/grpc"
	"auth/internal/services/auth"
	"auth/internal/services/manageAccount"
	"auth/internal/storage/postgres"
	"go.uber.org/zap"
	"time"
)

type App struct {
	GRPCServer *grpcApp.GRPC
	Storage    *postgres.Storage
}

func New(log *zap.Logger,
	gRPCPort int,
	port int,
	name, user, password string,
	tokenTTL time.Duration,
) *App {
	storage := postgres.New(port, name, user, password)
	authServ := auth.New(log, storage, tokenTTL)
	accountManagerSrv := managerAccount.New(log, storage)
	GRPCApp := grpcApp.New(log, authServ, accountManagerSrv, gRPCPort)

	return &App{
		GRPCServer: GRPCApp,
	}
}
