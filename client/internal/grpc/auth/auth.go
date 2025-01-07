// Пакет grpcAuth предоставляет клиент для взаимодействия с gRPC-сервисом аутентификации.
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

// New создает новый клиент AuthClient для подключения к gRPC-серверу аутентификации.
//
// Параметры:
//   - address: строка, содержащая адрес сервера.
//
// Возвращает:
//   - *AuthClient: клиент для взаимодействия с gRPC-сервисом.
//
// В случае невозможности установить соединение вызывается panic().
func New(address string) *AuthClient {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic("Failed to connect to Auth Server: " + err.Error())
	}
	return &AuthClient{client: auth.NewAuthClient(client)}
}

// Register регистрирует нового пользователя.
//
// Параметры:
//   - email: строка, содержащая email пользователя.
//   - password: строка, содержащая пароль пользователя.
//
// Возвращает:
//   - userID: уникальный идентификатор нового пользователя.
//   - err: ошибка в случае неудачной регистрации.
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

// Login выполняет вход пользователя и возвращает JWT-токен.
//
// Параметры:
//   - email: строка, содержащая email пользователя.
//   - password: строка, содержащая пароль пользователя.
//
// Возвращает:
//   - token: строка, содержащая JWT-токен.
//   - err: ошибка в случае неудачного логина.
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
