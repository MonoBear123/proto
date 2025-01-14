package jwt

import (
	"auth/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, duration time.Duration) (token string, err error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)

	claims["uid"] = user.ID
	claims["exp"] = time.Now().Add(duration).Unix() // Время истечения токена.
	claims["email"] = user.Email

	tokenString, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func NewManageAccountToken(id int64, duration time.Duration) (token string, err error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	claims["uid"] = id
	claims["exp"] = time.Now().Add(duration).Unix()
	tokenString, err := jwtToken.SignedString([]byte("activate"))
	if err != nil {
		return "", err

	}
	return tokenString, nil
}
