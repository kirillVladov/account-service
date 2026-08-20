package account_repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	row, err := db.Query(ctx, "SELECT id, email, password_hash, organization_id FROM account WHERE id = $1 AND organization_id = $2", id, organizationID)
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

func (r *Repository) Create(ctx context.Context, in dto.AccountCreateRequest) (dto.Account, error) {
	db := txManager.ExecutorFromContext(ctx, r.db)

	const query = `
		INSERT INTO account(
			email,
			password_hash,
			organization_id,
			created_at,
			updated_at
		) VALUES(
			@email,
			@password_hash,
			@organization_id,
			now(),
			now()
		)
		RETURNING
			id,
			email,
			password_hash,
			organization_id
	`

	args := pgx.NamedArgs{
		"email":           in.Email,
		"password_hash":   in.Password,
		"organization_id": in.OrganizationID,
	}

	rows, err := db.Query(ctx, query, args)
	if err != nil {
		return dto.Account{}, fmt.Errorf("exec query: insert account: %w", err)
	}

	raw, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[account])
	if err != nil {
		return dto.Account{}, fmt.Errorf("parse account: %w", err)
	}

	return convertToApplication(raw), nil
}
