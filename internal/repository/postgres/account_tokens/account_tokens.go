package account_tokens

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kirillVladov/account-service/internal/application/dto"
	"github.com/kirillVladov/account-service/internal/application/dto/errs"
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

func (r *Repository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, organizationID int64, tokenHash string, expiresAt time.Time) error {
	db := tx_manager.ExecutorFromContext(ctx, r.db)

	const query = `
		INSERT INTO auth_tokens (
			user_id,
			organization_id,
			token_type,
			token_hash,
			expires_at,
			created_at,
			updated_at
		) VALUES (
			@user_id,
			@organization_id,
			@token_type,
			@token_hash,
			@expires_at,
			NOW(),
			NOW()
		)
	`

	args := pgx.NamedArgs{
		"user_id":         userID,
		"organization_id": organizationID,
		"token_type":      string(dto.TokenTypeRefresh),
		"token_hash":      tokenHash,
		"expires_at":      expiresAt,
	}

	if _, err := db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("upsert token: %w", err)
	}

	return nil
}

func (r *Repository) CreateConfirmationToken(ctx context.Context, userID uuid.UUID, organizationID int64, tokenHash string, expiresAt time.Time) error {
	db := tx_manager.ExecutorFromContext(ctx, r.db)

	const query = `
		INSERT INTO auth_tokens (
			user_id,
			organization_id,
			token_type,
			token_hash,
			expires_at,
			created_at,
			updated_at
		) VALUES (
			@user_id,
			@organization_id,
			@token_type,
			@token_hash,
			@expires_at,
			NOW(),
			NOW()
		)
	`

	args := pgx.NamedArgs{
		"user_id":         userID,
		"organization_id": organizationID,
		"token_type":      string(dto.TokenTypeEmailConfirm),
		"token_hash":      tokenHash,
		"expires_at":      expiresAt,
	}

	if _, err := db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("create confirmation token: %w", err)
	}

	return nil
}

func (r *Repository) GetTokenByUserID(ctx context.Context, userID uuid.UUID, organizationID int64) (dto.AccountToken, error) {
	db := tx_manager.ExecutorFromContext(ctx, r.db)

	const query = `
		SELECT
			id,
			user_id,
			organization_id,
			token_type,
			token_hash,
			expires_at,
			revoked
		FROM auth_tokens 
		WHERE user_id = @user_id AND organization_id = @organization_id AND expires_at >= NOW() AND revoked = FALSE`

	args := pgx.NamedArgs{
		"user_id":         userID,
		"organization_id": organizationID,
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

func (r *Repository) GetByTokenHash(ctx context.Context, tokenHash string) (dto.AccountToken, error) {
	db := tx_manager.ExecutorFromContext(ctx, r.db)

	const query = `
		SELECT
			id,
			user_id,
			organization_id,
			token_type,
			token_hash,
			expires_at,
			revoked
		FROM auth_tokens
		WHERE token_hash = $1 AND expires_at >= NOW() AND revoked = FALSE
	`

	row, err := db.Query(ctx, query, tokenHash)
	if err != nil {
		return dto.AccountToken{}, fmt.Errorf("query token by hash: %w", err)
	}

	defer row.Close()

	raw, err := pgx.CollectOneRow(row, pgx.RowToStructByName[accountToken])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.AccountToken{}, errs.ErrTokenNotValid
		}

		return dto.AccountToken{}, fmt.Errorf("collect token row: %w", err)
	}

	return convertToApplication(raw), nil
}

func (r *Repository) DeactivateByUser(ctx context.Context, userID uuid.UUID, organizationID int64) error {
	db := tx_manager.ExecutorFromContext(ctx, r.db)

	const query = `
		UPDATE auth_tokens 
			SET
				revoked = TRUE
		WHERE user_id = @user_id AND organization_id = @organization_id
	`

	args := pgx.NamedArgs{
		"user_id":         userID,
		"organization_id": organizationID,
	}

	if _, err := db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("deactivate token: %w", err)
	}

	return nil
}
