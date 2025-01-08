// Пакет для создания JWT токенов.
package jwt

import (
	"auth/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// NewToken - создает новый JWT токен для пользователя.
//
// Параметры:
//   - user: Структура модели пользователя, для которого создается токен.
//   - duration: Время жизни токена.
//
// Возвращает:
//   - Строку с JWT токеном.
//   - Ошибку, если не удалось создать токен.
func NewToken(user models.User, duration time.Duration) (token string, err error) {
	// Создание нового JWT токена с алгоритмом HS256.
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	// Приведение claims токена к типу jwt.MapClaims.
	// jwt.MapClaims представляет собой  словарь (map),
	// где ключи — это строки, а значения могут быть любыми.
	claims := jwtToken.Claims.(jwt.MapClaims)

	claims["uid"] = user.ID
	claims["exp"] = time.Now().Add(duration).Unix() // Время истечения токена.
	claims["email"] = user.Email

	// SignedString принимает ключ,
	// который используется для создания подписи токена.
	// Подпись подтверждает, что токен был создан сервером и не был изменён после создания.
	tokenString, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
