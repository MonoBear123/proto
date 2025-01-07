// Пакет auth реализует логику аутентификации через gRPC.
// Включает обработку запросов на вход в систему (Login) и регистрацию нового пользователя (Register).
package auth

import (
	"context"
	"github.com/MonoBear123/proto/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Интерфейс AuthService описывает основные операции, связанные с аутентификацией,
type AuthService interface {
	// Login выполняет вход пользователя в систему.
	//
	// Принимает:
	// 	- ctx: контекст для управления сроками жизни запроса.
	// 	- username: адрес электронной почты пользователя.
	// 	- password: пароль пользователя.
	//
	// Возвращает:
	// 	- token: строку токена доступа, если вход успешен.
	// 	- err: ошибку, если вход не удался (например, неверные данные).
	Login(ctx context.Context, username, password string) (token string, err error)

	// RegisterNewUser регистрирует нового пользователя.
	//
	// Принимает:
	// 	- ctx: контекст для управления сроками жизни запроса.
	// 	- username: адрес электронной почты пользователя.
	// 	- password: пароль пользователя.
	//
	// Возвращает:
	// 	- userid: идентификатор зарегистрированного пользователя.
	// 	- err: ошибку, если регистрация не удалась (например, проблемы с сохранением данных).
	RegisterNewUser(ctx context.Context, username, password string) (userid int64, err error)
}

// serverAPI представляет собой структуру, которая реализует интерфейс AuthServer.
// Она содержит зависимость от AuthService, которая используется для выполнения логики аутентификации.
type serverAPI struct {
	auth.UnimplementedAuthServer // Встраивание автоматически сгенерированных методов.
	authSrv                      AuthService
}

// Register регистрирует сервер аутентификации на переданном gRPC сервере.
// Принимает:
// - gRPC: gRPC сервер, на который будет зарегистрирован сервис аутентификации.
// - authSrv: экземпляр сервиса аутентификации, который реализует бизнес-логику.
// Связывает gRPC сервер с сервисом.
func Register(gRPC *grpc.Server, authSrv AuthService) {
	auth.RegisterAuthServer(gRPC, &serverAPI{authSrv: authSrv})
}

// Login обрабатывает запросы на вход в систему.
// Принимает:
// - ctx: контекст для управления временем жизни запроса.
// - in: структуру запроса типа LoginRequest, содержащую параметры (email, password).
// Возвращает:
// - LoginResponse: ответ с токеном, если вход успешен.
// - ошибку, если вход не удался (например, неверные параметры или проблемы с сервером).
func (s *serverAPI) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
	if in.GetEmail() == "" || in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing params")
	}

	// Вызов сервиса аутентификации
	token, err := s.authSrv.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.LoginResponse{Token: token}, nil
}

// Register обрабатывает запросы на регистрацию нового пользователя.
// Принимает:
// - ctx: контекст для управления временем жизни запроса.
// - in: структуру запроса типа RegisterRequest, содержащую параметры (email, password).
// Возвращает:
// - RegisterResponse: ответ с идентификатором пользователя, если регистрация успешна.
// - ошибку, если регистрация не удалась (например, проблемы с сохранением данных).
func (s *serverAPI) Register(ctx context.Context, in *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if in.GetEmail() == "" || in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing params")
	}

	// Вызов сервиса для регистрации нового пользователя
	userID, err := s.authSrv.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &auth.RegisterResponse{
		UserId: userID,
	}, nil
}
