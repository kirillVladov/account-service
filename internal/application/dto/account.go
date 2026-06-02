package dto

import "github.com/google/uuid"

type UserRole string

const (
	UserRoleUser UserRole = "USER"
)

type Account struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Phone        string
	PasswordHash string
	TelegramID   string
}
