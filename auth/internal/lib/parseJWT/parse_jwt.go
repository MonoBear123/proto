package parseJWT

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func ParseToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("activate"), nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	expInt, ok := claims["exp"].(float64)
	if !ok {
		return 0, errors.New("invalid exp")
	}

	if time.Now().After(time.Unix(int64(expInt), 0)) {
		return 0, errors.New("expired")
	}
	userID, ok := claims["uid"].(int64)
	if !ok {
		return 0, errors.New("invalid userID")
	}
	return userID, nil
}
