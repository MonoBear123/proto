// Пакет authHandler предоставляет обработчики (handlers) для работы с аутентификацией и регистрацией пользователей.
// Использует gRPC-клиент для взаимодействия с соответствующим Auth-сервисом.
package authHandler

import (
	"client/internal/grpc/auth"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AuthHandler struct {
	client *grpcAuth.AuthClient
}

// New создает новый экземпляр AuthHandler.
//
// Параметры:
//   - client: указатель на gRPC-клиент AuthClient.
//
// Возвращает:
//   - *AuthHandler: экземпляр обработчика.
func New(client *grpcAuth.AuthClient) *AuthHandler {
	return &AuthHandler{client: client}
}

// Login обрабатывает запрос на вход пользователя.
//
// Принимает:
//   - ctx: контекст запроса Echo, из которого извлекаются данные формы (email и password).
//
// Возвращает:
//   - error: ошибку, если она произошла.
func (a *AuthHandler) Login(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	token, err := a.client.Login(email, password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}
	ctx.SetCookie(cookie)
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Token set successfully"})
}

// Register обрабатывает запрос на регистрацию нового пользователя.
//
// Принимает:
//   - ctx: контекст запроса Echo, из которого извлекаются данные формы (email и password).
//
// Возвращает:
//   - error: ошибку, если она произошла.
func (a *AuthHandler) Register(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	uid, err := a.client.Register(email, password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"uid": string(uid)})
}
