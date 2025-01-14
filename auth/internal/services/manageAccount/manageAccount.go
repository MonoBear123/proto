package managerAccount

import (
	"auth/internal/domain/models"
	"auth/internal/lib/jwt"
	"auth/internal/lib/parseJWT"
	"auth/internal/lib/sendMail"
	"context"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Storage interface {
	UpdatePassword(
		ctx context.Context,
		userID int64,
		password []byte,
	) error
	ActiveAccount(
		ctx context.Context,
		userID int64,
	) error
	ProvideUser(
		ctx context.Context,
		email string,
	) (user models.User, err error)
}

type AccountManager struct {
	logger  *zap.Logger
	storage Storage
}

func (am *AccountManager) ResetPassword(ctx context.Context, token, password string) (err error) {
	const op = "manageAccount.ResetPassword"
	log := am.logger.With(
		zap.String("op", op))
	userId, err := parseJWT.ParseToken(token)
	if err != nil {
		log.Error("failed to parse token", zap.Error(err))

		return err
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", zap.Error(err))

		return err
	}

	err = am.storage.UpdatePassword(ctx, userId, passHash)
	if err != nil {
		log.Error("failed to update password", zap.Error(err))

		return err
	}
	return nil
}

func (am *AccountManager) ActiveAccount(ctx context.Context, token string) (err error) {
	const op = "manageAccount.ActiveAccount"
	log := am.logger.With(zap.String("op", op))
	userId, err := parseJWT.ParseToken(token)
	if err != nil {
		log.Error("failed to parse token", zap.Error(err))

		return err
	}
	err = am.storage.ActiveAccount(ctx, userId)
	if err != nil {
		log.Error("failed to active account", zap.Error(err))

		return err
	}
	return nil
}

func (am *AccountManager) ForgotPassword(ctx context.Context, email string) (err error) {
	const op = "manageAccount.ForgotPassword"
	log := am.logger.With(
		zap.String("op", op))
	user, err := am.storage.ProvideUser(ctx, email)
	if err != nil {
		log.Error("failed to provide user", zap.Error(err))

		return err

	}
	TokenTTL := time.Hour * 72
	token, err := jwt.NewManageAccountToken(user.ID, TokenTTL)
	if err != nil {
		log.Error("failed to generate token", zap.Error(err))

		return err
	}

	err = sendMail.SendMessagee(email, "Сброс пароля", token)
	if err != nil {
		log.Error("failed to send email", zap.Error(err))

		return err
	}
	return nil
}

func New(logger *zap.Logger, storage Storage) *AccountManager {
	return &AccountManager{
		logger:  logger,
		storage: storage,
	}
}
