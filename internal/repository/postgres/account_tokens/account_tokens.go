package account_tokens

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kirillVladov/account-service/internal/application/dto"
	tx_manager "github.com/kirillVladov/account-service/pkg/tx"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	db := tx_manager.ExecutorFromContext(ctx, r.db)

	const query = `
		INSERT INTO auth_tokens (
			user_id,   
			token_hash,
			expires_at,
			created_at,
			updated_at
		) VALUES (
			@user_id,
			@token_hash,
			@expires_at,
			NOW(),
			NOW()
		)
	`

	args := pgx.NamedArgs{
		"user_id":    userID,
		"token_hash": tokenHash,
		"expires_at": expiresAt,
	}

	if err := db.QueryRow(ctx, query, args); err != nil {
		return fmt.Errorf("upsert token: %w", err)
	}

	return nil
}

func (r *Repository) GetTokenByUserID(ctx context.Context, userID uuid.UUID) (dto.AccountToken, error) {
	db := tx_manager.ExecutorFromContext(ctx, r.db)

	const query = `
		SELECT
			id,
			user_id,
			token_type,
			token_hash,
			expires_at,
			revoked
		FROM auth_tokens 
		WHERE user_id = @user_id AND expires_at >= NOW() AND revoked = FALSE`

	args := pgx.NamedArgs{
		"user_id": userID,
	}

	rows, err := db.Query(ctx, query, args)
	if err != nil {
		return dto.AccountToken{}, fmt.Errorf("getting tokens: %w", err)
	}

	rawAccount, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[accountToken])
	if err != nil {
		return dto.AccountToken{}, fmt.Errorf("parse auth_tokens: %w", err)
	}

	return convertToApplication(rawAccount), nil
}

func (r *Repository) DeactivateByUser(ctx context.Context, userID uuid.UUID) error {
	db := tx_manager.ExecutorFromContext(ctx, r.db)

	const query = `
		UPDATE auth_tokens 
			SET
				revoked = TRUE
		WHERE user_id = @user_id
	`

	args := pgx.NamedArgs{
		"user_id": userID,
	}

	if _, err := db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("deactivate token: %w", err)
	}

	return nil
}
