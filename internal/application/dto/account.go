package dto

import (
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeRefresh TokenType = "refresh"
	TokenTypeAccess  TokenType = "access"
)

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

type AccountToken struct {
	ID        int64
	UserID    uuid.UUID
	TokenType TokenType
	TokenHash string
	ExpiresAt time.Time
	Revoked   bool
}
