package managerAccount

import (
	"context"
	"github.com/MonoBear123/proto/protos/gen/go/auth"
	"google.golang.org/grpc"
)

type AccountManagerService interface {
	ResetPassword(ctx context.Context, email, password string) (err error)
	ActiveAccount(ctx context.Context, token string) (err error)
	ForgotPassword(ctx context.Context, email string) (err error)
}
type serverAPI struct {
	auth.UnimplementedAuthServer
	accountManagerSrv AccountManagerService
}

func Register(gRPC *grpc.Server, accountManagerSrv AccountManagerService) {
	auth.RegisterAuthServer(gRPC, &serverAPI{accountManagerSrv: accountManagerSrv})
}

func (s *serverAPI) ForgotPassword(ctx context.Context, in *auth.ForgotPasswordRequest) (*auth.ForgotPasswordResponse, error) {
	email := in.GetEmail()
	err := s.accountManagerSrv.ForgotPassword(ctx, email)
	if err != nil {
		return nil, err
	}
	return &auth.ForgotPasswordResponse{
		Message: "email send",
	}, nil
}

func (s *serverAPI) ResetPassword(ctx context.Context, in *auth.ResetPasswordRequest) (*auth.ResetPasswordResponse, error) {
	password := in.GetPassword()
	token := in.GetToken()
	err := s.accountManagerSrv.ResetPassword(ctx, token, password)
	if err != nil {

		return nil, err
	}
	return &auth.ResetPasswordResponse{
		Message: "password reset",
	}, nil
}

func (s *serverAPI) ActiveAccount(ctx context.Context, in *auth.ActiveAccountRequest) (*auth.ActiveAccountResponse, error) {
	token := in.GetToken()
	err := s.accountManagerSrv.ActiveAccount(ctx, token)
	if err != nil {
		return nil, err
	}
	return &auth.ActiveAccountResponse{Message: "account activated"}, nil

}
