package account_tokens

import (
	"time"

	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
)

type accountToken struct {
	ID        int64     `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	TokenType string    `db:"token_type"`
	TokenHash string    `db:"token_hash"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
}

func convertToApplication(in accountToken) dto.AccountToken {
	return dto.AccountToken{
		ID:        in.ID,
		UserID:    in.UserID,
		TokenType: dto.TokenType(in.TokenHash),
		TokenHash: in.TokenHash,
		ExpiresAt: in.ExpiresAt,
		Revoked:   in.Revoked,
	}
}
