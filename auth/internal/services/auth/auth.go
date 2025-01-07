// Пакет auth реализует сервис аутентификации и регистрации пользователей.
// Включает обработку логина, регистрацию нового пользователя, генерацию токенов
// и взаимодействие с слоем работы с данными.
package auth

import (
	"auth/internal/domain/models"
	"auth/internal/lib/jwt"
	"auth/internal/storage"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Storage - интерфейс для работы с хранилищем пользователей.
// Определяет методы для создания, обновления и получения пользователей.
type Storage interface {
	// UpdateUser обновляет данные пользователя в хранилище.
	// Принимает:
	// 	- ctx: контекст для управления временем жизни запроса.
	// 	- user: обновленные данные пользователя.
	// Возвращает:
	// 	- err: ошибку, при возникновении.
	UpdateUser(
		ctx context.Context,
		user models.User,
	) (err error)

	// SaveUser сохраняет нового пользователя в хранилище.
	// Принимает:
	// 	- ctx: контекст для управления временем жизни запроса.
	// 	- email: электронная почта пользователя.
	// 	- passHash: хэш пароля пользователя.
	// Возвращает:
	// 	- uid: уникальный идентификатор нового пользователя.
	// 	- err: ошибку, при возникновении.
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)

	// ProvideUser извлекает данные пользователя из хранилища по email.
	// Принимает:
	// - ctx: контекст для управления временем жизни запроса.
	// - email: электронная почта пользователя.
	// Возвращает:
	// - user: данные пользователя.
	// - err: ошибку, при возникновении.
	ProvideUser(
		ctx context.Context,
		email string,
	) (user models.User, err error)
}
type Auth struct {
	logger   *zap.Logger
	storage  Storage
	tokenTTL time.Duration
}

// New - конструктор для создания нового экземпляра Auth.
// Принимает:
//   - logger: логгер для ведения логов.
//   - storage: слой работы с хранилищем пользователей.
//   - tokenTTL: время жизни токена.
//
// Возвращает:
//   - Указатель на структуру Auth.
func New(logger *zap.Logger, storage Storage, tokenTTL time.Duration) *Auth {
	return &Auth{
		logger:   logger,
		storage:  storage,
		tokenTTL: tokenTTL,
	}
}

// Login выполняет аутентификацию пользователя.
// Принимает:
//   - ctx: контекст для управления временем жизни запроса.
//   - email: электронная почта пользователя.
//   - password: пароль пользователя.
//
// Возвращает:
//   - token: JWT токен, если аутентификация прошла успешно.
//   - err: ошибку, если аутентификация не удалась.
func (a *Auth) Login(ctx context.Context, email, password string) (token string, err error) {
	const op = "auth.Login"
	log := a.logger.With(
		zap.String("op", op))

	log.Info("login user")
	// Работа сервисного слоя со слоем данных.
	// Получение данных пользователя из хранилища по-указанному email.
	user, err := a.storage.ProvideUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.logger.Warn("user not found", zap.Error(err))

			return "", fmt.Errorf("invalid data")
		}

		a.logger.Error("failed to get user", zap.Error(err))

		return "", fmt.Errorf("internal server error")
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.logger.Warn("invalid date", zap.Error(err))

		return "", fmt.Errorf("invalid data")
	}
	token, err = jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		a.logger.Error("failed to generate token", zap.Error(err))

		return "", err
	}
	return token, nil

}

// RegisterNewUser регистрирует нового пользователя в системе.
// Принимает:
//   - ctx: контекст для управления временем жизни запроса.
//   - email: электронная почта нового пользователя.
//   - password: пароль нового пользователя.
//
// Возвращает:
//   - userid: уникальный идентификатор нового пользователя, если регистрация прошла успешно.
//   - err: ошибку, если регистрация не удалась.
func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (userid int64, err error) {
	const op = "auth.RegisterNewUser"
	log := a.logger.With(
		zap.String("op", op),
	)
	log.Info("register new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", zap.Error(err))

		return 0, err
	}
	// Работа сервисного слоя со слоем данных.
	// Сохранение данных пользователя в хранилище.
	id, err := a.storage.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.logger.Warn("user already exists", zap.Error(err))

			return 0, nil
		}
		log.Error("failed to save user", zap.Error(err))

		return 0, fmt.Errorf("failed to registration")
	}
	return id, nil

}
