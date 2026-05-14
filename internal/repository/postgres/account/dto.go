package account_repo

import (
	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
)

type account struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	Name         string    `db:"name"`
	Token        string    `db:"token"`
	RefreshToken string    `db:"refresh_token"`
}

func convertToApplication(in account) dto.Account {
	return dto.Account{
		ID:           in.ID,
		Email:        in.Email,
		Name:         in.Name,
		Token:        in.Token,
		RefreshToken: in.RefreshToken,
	}
}

func convertToRepository(in dto.Account) account {
	return account{
		ID:           in.ID,
		Email:        in.Email,
		Name:         in.Name,
		Token:        in.Token,
		RefreshToken: in.RefreshToken,
	}
}
