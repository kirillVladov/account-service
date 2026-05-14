package dto

import "github.com/google/uuid"

type Account struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Token        string
	RefreshToken string
}
