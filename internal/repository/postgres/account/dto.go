package account_repo

import (
	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
)

type account struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	Name         string    `db:"name"`
	PasswordHash string    `db:"password_hash"`
	Phone        string    `db:"phone"`
	TelegramID   string    `db:"telegram_id"`
}

func convertToApplication(in account) dto.Account {
	return dto.Account{
		ID:           in.ID,
		Email:        in.Email,
		Name:         in.Name,
		PasswordHash: in.PasswordHash,
		Phone:        in.Phone,
		TelegramID:   in.TelegramID,
	}
}

func convertToRepository(in dto.Account) account {
	return account{
		ID:           in.ID,
		Email:        in.Email,
		Name:         in.Name,
		PasswordHash: in.PasswordHash,
		Phone:        in.Phone,
		TelegramID:   in.TelegramID,
	}
}
