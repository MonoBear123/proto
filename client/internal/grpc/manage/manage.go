package manageGRPC

import (
	"context"
	manage "github.com/MonoBear123/proto/protos/gen/go/auth"
	"google.golang.org/grpc"
	"time"
)

type ManagerClient struct {
	client manage.AccountManagerClient
}

func New(address string) *ManagerClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic("Failed to connect to AccountManager Server: " + err.Error())
	}
	return &ManagerClient{client: manage.NewAccountManagerClient(client)}
}

func (a *ManagerClient) ForgotPass(email string) error {
	_, err := a.client.ForgotPassword(context.Background(), &manage.ForgotPasswordRequest{
		Email: email,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *ManagerClient) ActivateAccount(token string) error {
	_, err := a.client.ActiveAccount(context.Background(), &manage.ActiveAccountRequest{
		Token: token,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *ManagerClient) ResetPassword(token, password string) error {
	_, err := a.client.ResetPasword(context.Background(), &manage.ResetPasswordRequest{
		Token:    token,
		Password: password,
	})
	if err != nil {
		return err
	}
	return nil
}
