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

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (dto.Account, error) {
	db := txManager.ExecutorFromContext(ctx, r.db)

	row, err := db.Query(ctx, "SELECT id, email, name, telegram_id, password_hash, phone FROM account WHERE id = $1", id)
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

func (r *Repository) Create(ctx context.Context, account dto.Account) error {
	db := txManager.ExecutorFromContext(ctx, r.db)

	const query = `
		INSERT INTO account(
			id,
			email,
			name,
			password_hash,
			telegram_id,
			phone,
			created_at,
			updated_at
		) VALUES(
			@id,
			@email,
			@name,
			@password_hash,
			@telegram_id,
			@phone,
			now(),
			now()
		)
	`

	rec := convertToRepository(account)

	args := pgx.NamedArgs{
		"id":            rec.ID,
		"email":         rec.Email,
		"phone":         rec.Phone,
		"name":          rec.Name,
		"password_hash": rec.PasswordHash,
		"telegram_id":   rec.TelegramID,
	}

	_, err := db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("exec query: insert account: %w", err)
	}

	return nil
}

// func (r *Repository) UpdateToken(ctx context.Context, id uuid.UUID, token string) error {
// 	db := txManager.ExecutorFromContext(ctx, r.db)

// 	_, err := db.Exec(ctx, "INSERT INTO auth_tokens(user_id, token_hash, expires_at) VALUES($1, $2, NOW() + INTERVAL '7 days')", id, token)
// 	if err != nil {
// 		return fmt.Errorf("exec query: update token: %w", err)
// 	}

// 	return nil
// }

func (r *Repository) GetByTelegramID(ctx context.Context, telegramID string) (dto.Account, error) {
	db := txManager.ExecutorFromContext(ctx, r.db)

	row, err := db.Query(ctx, "SELECT id, email, name, telegram_id, password_hash, phone FROM account WHERE telegram_id = $1", telegramID)
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
