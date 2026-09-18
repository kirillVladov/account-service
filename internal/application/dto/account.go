package dto

import (
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeRefresh          TokenType = "refresh_token"
	TokenTypeEmailConfirm     TokenType = "email_confirmation_token"
)

type UserRole string

const (
	UserRoleUser UserRole = "USER"
)

type Account struct {
	ID             uuid.UUID
	Email          string
	PasswordHash   string
	OrganizationID int64
	IsConfirmed    bool
	IsBlocked      bool
}

type AccountCreateRequest struct {
	Email          string
	Password       string
	OrganizationID int64
	IdempotencyKey uuid.UUID
}

type AccountToken struct {
	ID             int64
	UserID         uuid.UUID
	OrganizationID int64
	TokenType      TokenType
	TokenHash      string
	ExpiresAt      time.Time
	Revoked        bool
}
