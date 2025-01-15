package manageHandler

import (
	"client/internal/grpc/manage"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ManageHandler struct {
	client *manageGRPC.ManagerClient
}

func New(client *manageGRPC.ManagerClient) *ManageHandler {
	return &ManageHandler{client: client}
}

func (a *ManageHandler) ForgotPass(ctx echo.Context) error {
	email := ctx.FormValue("email")

	err := a.client.ForgotPass(email)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Пользователь не найден", err)})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Письмо отправлено на почту"})

}

func (a *ManageHandler) ResetPass(ctx echo.Context) error {
	password := ctx.FormValue("password")
	token := ctx.QueryParam("token")
	err := a.client.ResetPassword(token, password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Пароль изменен"})
}

func (a *ManageHandler) ActivateAccount(ctx echo.Context) error {
	token := ctx.QueryParam("token")
	err := a.client.ActivateAccount(token)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Аккаунт активирован"})

}
