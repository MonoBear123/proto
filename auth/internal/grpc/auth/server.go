package auth

import (
	"context"
	"github.com/MonoBear123/proto/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (token string, err error)
	RegisterNewUser(ctx context.Context, username, password string) (userid int64, err error)
}

type serverAPI struct {
	auth.UnimplementedAuthServer // Встраивание автоматически сгенерированных методов.
	authSrv                      AuthService
}

func Register(gRPC *grpc.Server, authSrv AuthService) {
	auth.RegisterAuthServer(gRPC, &serverAPI{authSrv: authSrv})
}

func (s *serverAPI) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	if in.GetEmail() == "" || in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing params")
	}

	token, err := s.authSrv.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, in *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if in.GetEmail() == "" || in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing params")
	}

	userID, err := s.authSrv.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.RegisterResponse{
		UserId: userID,
	}, nil
}
