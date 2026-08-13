package account_repo

import (
	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
)

type account struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
}

func convertToApplication(in account) dto.Account {
	return dto.Account{
		ID:           in.ID,
		Email:        in.Email,
		PasswordHash: in.PasswordHash,
	}
}

func convertToRepository(in dto.Account) account {
	return account{
		ID:           in.ID,
		Email:        in.Email,
		PasswordHash: in.PasswordHash,
	}
}
