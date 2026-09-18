package account_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kirillVladov/account-service/internal/application/dto"
	"github.com/kirillVladov/account-service/internal/application/dto/errs"
	txManager "github.com/kirillVladov/account-service/pkg/tx"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID, organizationID int64) (dto.Account, error) {
	db := txManager.ExecutorFromContext(ctx, r.db)

	row, err := db.Query(ctx, "SELECT id, email, password_hash, organization_id, confirmed, blocked FROM account WHERE id = $1 AND organization_id = $2", id, organizationID)
	if err != nil {
		return dto.Account{}, fmt.Errorf("query account: %w", err)
	}

	defer row.Close()

	account, err := pgx.CollectOneRow(row, pgx.RowToStructByName[account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.Account{}, errs.ErrAccountNotFound
		}

		return dto.Account{}, fmt.Errorf("collect account row: %w", err)
	}

	return convertToApplication(account), nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string, organizationID int64) (dto.Account, error) {
	db := txManager.ExecutorFromContext(ctx, r.db)

	const query = `
		SELECT
			id,
			email,
			password_hash,
			organization_id,
			confirmed,
			blocked
		FROM account
		WHERE email = $1 AND organization_id = $2
	`

	row, err := db.Query(ctx, query, email, organizationID)
	if err != nil {
		return dto.Account{}, fmt.Errorf("query account by email: %w", err)
	}

	account, err := pgx.CollectOneRow(row, pgx.RowToStructByName[account])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.Account{}, errs.ErrAccountNotFound
		}

		return dto.Account{}, fmt.Errorf("collect account row: %w", err)
	}

	return convertToApplication(account), nil
}

func (r *Repository) Create(ctx context.Context, in dto.AccountCreateRequest) (dto.Account, error) {
	db := txManager.ExecutorFromContext(ctx, r.db)

	const query = `
		INSERT INTO account(
			email,
			password_hash,
			organization_id,
			idempotency_key,
			created_at,
			updated_at
		) VALUES(
			@email,
			@password_hash,
			@organization_id,
			@idempotency_key,
			now(),
			now()
		)
		RETURNING
			id,
			email,
			password_hash,
			organization_id,
			confirmed,
			blocked
	`

	args := pgx.NamedArgs{
		"email":           in.Email,
		"password_hash":   in.Password,
		"organization_id": in.OrganizationID,
		"idempotency_key": in.IdempotencyKey,
	}

	rows, err := db.Query(ctx, query, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return dto.Account{}, errs.ErrAccountAlreadyExists
		}

		return dto.Account{}, fmt.Errorf("exec query: insert account: %w", err)
	}

	raw, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account])
	if err != nil {
		return dto.Account{}, fmt.Errorf("parse account: %w", err)
	}

	return convertToApplication(raw), nil
}

func (r *Repository) SetBlocked(ctx context.Context, id uuid.UUID, organizationID int64, blocked bool) error {
	db := txManager.ExecutorFromContext(ctx, r.db)

	const query = `
		UPDATE account
		SET blocked = $3, updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`

	if _, err := db.Exec(ctx, query, id, organizationID, blocked); err != nil {
		return fmt.Errorf("set blocked: %w", err)
	}

	return nil
}

func (r *Repository) SetConfirmed(ctx context.Context, id uuid.UUID, organizationID int64) error {
	db := txManager.ExecutorFromContext(ctx, r.db)

	const query = `
		UPDATE account
		SET confirmed = TRUE, updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`

	if _, err := db.Exec(ctx, query, id, organizationID); err != nil {
		return fmt.Errorf("set confirmed: %w", err)
	}

	return nil
}
