package auth

import (
	"auth/internal/domain/models"
	"auth/internal/lib/jwt"
	"auth/internal/lib/sendMail"
	"auth/internal/storage"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	TokenTLForActivate = time.Hour * 72
)

type Storage interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
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

func New(logger *zap.Logger, storage Storage, tokenTTL time.Duration) *Auth {
	return &Auth{
		logger:   logger,
		storage:  storage,
		tokenTTL: tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string) (token string, err error) {
	const op = "auth.Login"
	log := a.logger.With(
		zap.String("op", op))

	log.Info("login user")
	user, err := a.storage.ProvideUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.logger.Warn("user not found", zap.Error(err))

			return "", fmt.Errorf("invalid data")
		}

		a.logger.Error("failed to get user", zap.Error(err))

		return "", fmt.Errorf("internal server error")
	}
	if !user.ActiveAccount {
		return "", fmt.Errorf("invalid date")
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
	id, err := a.storage.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", zap.Error(err))

			return 0, nil
		}
		log.Error("failed to save user", zap.Error(err))

		return 0, fmt.Errorf("failed to registration")
	}
	token, err := jwt.NewManageAccountToken(id, TokenTLForActivate)
	err = sendMail.SendMessagee(email, "Активация аккаунта", token)
	if err != nil {
		log.Error("failed to send mail", zap.Error(err))

		return 0, err
	}
	return id, nil

}
