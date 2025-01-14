package grpcAuth

import (
	"context"
	"github.com/MonoBear123/proto/protos/gen/go/auth"
	"google.golang.org/grpc"
	"time"
)

type AuthClient struct {
	client auth.AuthClient
}

func New(address string) *AuthClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic("Failed to connect to Auth Server: " + err.Error())
	}
	return &AuthClient{client: auth.NewAuthClient(client)}
}

func (a *AuthClient) Register(email string, password string) (userID int64, err error) {
	// Создание контекста с таймаутом для попытки подключения.
	res, err := a.client.Register(context.Background(), &auth.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return -1, err
	}

	return res.UserId, nil
}

func (a *AuthClient) Login(email string, password string) (token string, err error) {
	res, err := a.client.Login(context.Background(), &auth.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err

	}
	return res.Token, nil
}

func (a *AuthClient) ForgotPass(email string) error {
	_, err := a.client.ForgotPassword(context.Background(), &auth.ForgotPasswordRequest{
		Email: email,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthClient) ActivateAccount(token string) error {
	_, err := a.client.ActiveAccount(context.Background(), &auth.ActiveAccountRequest{
		Token: token,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthClient) ResetPassword(token, password string) error {
	_, err := a.client.ResetPasword(context.Background(), &auth.ResetPasswordRequest{
		Token:    token,
		Password: password,
	})
	if err != nil {
		return err
	}
	return nil
}
