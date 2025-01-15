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

func (s *Storage) UpdatePassword(ctx context.Context, userID int64, password []byte) error {
	const op = "storage.postgres.UpdatePassword"
	_, err := s.db.ExecEx(ctx, "UPDATE users SET passHash = $1 WHERE id = $2;", nil, password, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ActiveAccount(ctx context.Context, userID int64) error {
	const op = "storage.postgres.ActiveAccount"
	query := "UPDATE users SET activateAccount = $1 WHERE id = $2;"
	_, err := s.db.ExecEx(ctx, query, nil, true, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

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

func (s *Storage) ProvideUser(ctx context.Context, email string) (user models.User, err error) {
	const op = "storage.postgres.ProvideUser"

	err = s.db.QueryRowEx(
		ctx,
		"SELECT id, email, passHash,activateAccount FROM users WHERE email = $1", nil,
		email,
	).Scan(&user.ID, &user.Email, &user.PassHash, &user.ActiveAccount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, fmt.Errorf("%s: user not found: %w", op, err)
		}
		return models.User{}, fmt.Errorf("%s: query failed: %w", op, err)
	}

	return user, nil
}

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
        passHash BYTEA NOT NULL,
        activateAccount BOOLEAN DEFAULT FALSE
    );
`
	_, err = conn.Exec(createTableSQL)
	if err != nil {
		panic(fmt.Errorf("Failed to create table:%w", err))
	}
	return &Storage{db: conn}

}

func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		panic("Failed to close database")
	}
}
