package account_repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kirillVladov/account-service/internal/application/dto"
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

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (dto.Account, error) {
	db := txManager.ExecutorFromContext(ctx, r.db)

	row, err := db.Query(ctx, "SELECT id, email, name, token, refresh_token FROM account WHERE id = $1", id)
	if err != nil {
		return dto.Account{}, fmt.Errorf("query account: %w", err)
	}

	defer row.Close()

	account, err := pgx.CollectOneRow(row, pgx.RowToStructByName[account])
	if err != nil {
		return dto.Account{}, fmt.Errorf("collect account row: %w", err)
	}

	return convertToApplication(account), nil
}

func (r *Repository) Create(ctx context.Context, account dto.Account) error {
	db := txManager.ExecutorFromContext(ctx, r.db)

	const query = `
		INSERT INTO account(
			id,
			email,
			name,
			token,
			refresh_token
		) VALUES(
			@id,
			@email,
			@name,
			@token,
			@refresh_token
		)
	`

	rec := convertToRepository(account)

	args := pgx.NamedArgs{
		"id":            rec.ID,
		"email":         rec.Email,
		"name":          rec.Name,
		"token":         rec.Token,
		"refresh_token": rec.RefreshToken,
	}

	row, err := db.Query(ctx, query, args)
	if err != nil {
		return fmt.Errorf("query account: %w", err)
	}

	defer row.Close()

	return nil
}

func (r *Repository) UpdateToken(ctx context.Context, id uuid.UUID, token string) error {
	db := txManager.ExecutorFromContext(ctx, r.db)

	_, err := db.Exec(ctx, "UPDATE account SET token = $1 WHERE id = $2", token, id)
	if err != nil {
		return fmt.Errorf("exec query: update token: %w", err)
	}

	return nil
}
