package authHandler

import (
	"client/internal/grpc/auth"
	"github.com/labstack/echo/v4"
	"net/http"
	"regexp"
)

type AuthHandler struct {
	client *grpcAuth.AuthClient
}

func New(client *grpcAuth.AuthClient) *AuthHandler {
	return &AuthHandler{client: client}
}
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

func (a *AuthHandler) Register(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	uid, err := a.client.Register(email, password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]int64{"uid": uid})
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
