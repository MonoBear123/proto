// Пакет parseJWT предоставляет функционал для парсинга и проверки JWT-токенов.
package parseJWT

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// ParseToken разбирает JWT-токен и проверяет его валидность.
//
// Параметры:
// - tokenString: строка, содержащая JWT-токен.
//
// Возвращает:
// - error: ошибку, если токен невалиден или истек.
func ParseToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid claims")
	}

	expInt, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid exp")
	}

	if time.Now().After(time.Unix(int64(expInt), 0)) {
		return errors.New("expired")
	}

	return nil
}
