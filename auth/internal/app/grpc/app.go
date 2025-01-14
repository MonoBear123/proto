package grpcApp

import (
	authgrpc "auth/internal/grpc/auth"
	"auth/internal/grpc/managerAccount"
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

func New(log *zap.Logger, authServ authgrpc.AuthService, accountManager managerAccount.AccountManagerService, port int) *GRPC {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer, authServ)
	managerAccount.Register(gRPCServer, accountManager)
	return &GRPC{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (g *GRPC) MustRun() {
	if err := g.Run(); err != nil {
		panic(err)
	}
}

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

func (g *GRPC) Stop() {
	const op = "grpcApp.Stop"
	log := g.log.With(zap.String("op", op))
	log.Info("stopping gRPC server")

	g.gRPCServer.GracefulStop()
}
