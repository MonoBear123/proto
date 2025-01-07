// Пакет storage определяет общие ошибки, используемые слоями данных и сервисами.
package storage

import "errors"

var (
	ErrUserExists   = errors.New("User already exists")
	ErrUserNotFound = errors.New("User not found")
)
