// Пакет postgres реализует взаимодействие с базой данных PostgreSQL для работы с пользователями.
package postgres

import (
	"auth/internal/domain/models"
	"auth/internal/storage"
	"context"
	"fmt"
	"github.com/jackc/pgx"
)

type Storage struct {
	db *pgx.Conn
}

func (s *Storage) UpdateUser(ctx context.Context, user models.User) (err error) {
	//TODO implement me
	return
}

// SaveUser сохраняет нового пользователя в хранилище.
// Принимает:
//   - ctx: контекст для управления запросом и его отмены.
//   - email: email пользователя (уникальное значение).
//   - passHash: хеш пароля пользователя.
//
// Возвращает:
//   - uid: ID созданного пользователя.
//   - err: ошибка, если сохранение пользователя не удалось.
//   - storage.ErrUserExists: пользователь с указанным email уже существует.
//   - Другие ошибки, связанные с запросом в базу данных.
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	const op = "storage.postgres.SaveUser"
	var existingUser int

	err = s.db.QueryRowEx(
		ctx,
		"SELECT id FROM users WHERE email = $1", nil, email,
	).Scan(&existingUser)
	if err == nil {
		return 0, fmt.Errorf("user already exists: %w", storage.ErrUserExists)
	} else if err != pgx.ErrNoRows {
		return 0, fmt.Errorf("error checking for existing user: %w", err)
	}
	err = s.db.QueryRowEx(
		ctx,
		"INSERT INTO users (email, passHash) VALUES ($1, $2) RETURNING id", nil,
		email, passHash,
	).Scan(&uid)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return uid, nil
}

// ProvideUser извлекает данные пользователя.
// Принимает:
//   - ctx: контекст для управления запросом и его отмены.
//   - email: email пользователя, данные которого нужно извлечь.
//
// Возвращает:
//   - user: структура модели пользователя, содержащая данные из хранилища.
//   - err: ошибка, если извлечение данных не удалось.
//   - pgx.ErrNoRows: пользователь с указанным email не найден.
//   - Другие ошибки, связанные с запросом в базу данных.
func (s *Storage) ProvideUser(ctx context.Context, email string) (user models.User, err error) {
	const op = "storage.postgres.ProvideUser"

	err = s.db.QueryRowEx(
		ctx,
		"SELECT id, email, passHash FROM users WHERE email = $1", nil,
		email,
	).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, fmt.Errorf("%s: user not found: %w", op, err)
		}
		return models.User{}, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return user, nil
}

// New создает новое подключение к базе данных PostgreSQL и инициализирует таблицу пользователей.
// Принимает:
//   - port: порт для подключения к базе данных.
//   - name: имя базы данных.
//   - user: имя пользователя для подключения.
//   - password: пароль для подключения.
//
// Возвращает:
//   - указатель на структуру Storage.
//
// В случае ошибки подключения или создания таблицы вызывает panic().
func New(port int, name, user, password string) *Storage {
	connConf := pgx.ConnConfig{
		Host:     "authDB",
		Port:     uint16(port),
		User:     user,
		Password: password,
		Database: name,
	}
	conn, err := pgx.Connect(connConf)
	if err != nil {
		panic(fmt.Errorf("Failed to connect to database:%w", err))
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        passHash BYTEA NOT NULL
    );
`
	_, err = conn.Exec(createTableSQL)
	if err != nil {
		panic(fmt.Errorf("Failed to create table:%w", err))
	}

	return &Storage{db: conn}

}

// Close закрывает подключение к базе данных.
// В случае ошибки подключения или создания таблицы вызывает panic().
func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		panic("Failed to close database")
	}
}
