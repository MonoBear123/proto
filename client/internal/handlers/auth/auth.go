// Пакет authHandler предоставляет обработчики (handlers) для работы с аутентификацией и регистрацией пользователей.
// Использует gRPC-клиент для взаимодействия с соответствующим Auth-сервисом.
package authHandler

import (
	"client/internal/grpc/auth"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
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

func ValidateDate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		email := ctx.FormValue("email")
		password := ctx.FormValue("password")
		if email == "" || password == "" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Email and password are required"})
		}

		emailRegex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
		emailMatch, _ := regexp.MatchString(emailRegex, email)
		if !emailMatch {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid email format"})
		}

		if len(password) < 8 {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Password must be at least 8 characters long"})
		}

		return next(ctx)
	}
}

func (a *AuthHandler) ForgotPass(ctx echo.Context) error {
	email := ctx.FormValue("email")

	err := a.client.ForgotPass(email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Пользователь не найден", err)})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Письмо отправлено на почту"})

}

func (a *AuthHandler) ResetPass(ctx echo.Context) error {
	password := ctx.FormValue("password")
	token := ctx.QueryParam("token")
	err := a.client.ResetPassword(token, password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Пароль изменен"})
}

func (a *AuthHandler) ActivateAccount(ctx echo.Context) error {
	token := ctx.QueryParam("token")
	err := a.client.ActivateAccount(token)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Аккаунт активирован"})

}
