package managerAccount

import (
	"context"
	manage "github.com/MonoBear123/proto/protos/gen/go/auth"
	"google.golang.org/grpc"
)

type AccountManagerService interface {
	ResetPassword(ctx context.Context, email, password string) (err error)
	ActiveAccount(ctx context.Context, token string) (err error)
	ForgotPassword(ctx context.Context, email string) (err error)
}
type serverAPI struct {
	manage.UnimplementedAccountManagerServer
	accountManagerSrv AccountManagerService
}

func Register(gRPC *grpc.Server, accountManagerSrv AccountManagerService) {
	manage.RegisterAccountManagerServer(gRPC, &serverAPI{accountManagerSrv: accountManagerSrv})
}

func (s *serverAPI) ForgotPassword(ctx context.Context, in *manage.ForgotPasswordRequest) (*manage.ForgotPasswordResponse, error) {
	email := in.GetEmail()
	err := s.accountManagerSrv.ForgotPassword(ctx, email)
	if err != nil {
		return nil, err
	}
	return &manage.ForgotPasswordResponse{
		Message: "email send",
	}, nil
}

func (s *serverAPI) ResetPasword(ctx context.Context, in *manage.ResetPasswordRequest) (*manage.ResetPasswordResponse, error) {
	password := in.GetPassword()
	token := in.GetToken()
	err := s.accountManagerSrv.ResetPassword(ctx, token, password)
	if err != nil {

		return nil, err
	}
	return &manage.ResetPasswordResponse{
		Message: "password reset",
	}, nil
}

func (s *serverAPI) ActiveAccount(ctx context.Context, in *manage.ActiveAccountRequest) (*manage.ActiveAccountResponse, error) {
	token := in.GetToken()
	err := s.accountManagerSrv.ActiveAccount(ctx, token)
	if err != nil {
		return nil, err
	}
	return &manage.ActiveAccountResponse{Message: "account activated"}, nil

}
