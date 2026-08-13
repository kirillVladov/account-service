package organization_repo

import (
	"context"
	"errors"
	"fmt"

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

func (r *Repository) GetByID(ctx context.Context, id int64) (dto.Organization, error) {
	db := txManager.ExecutorFromContext(ctx, r.db)

	rows, err := db.Query(ctx, "SELECT id, name FROM organizations WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return dto.Organization{}, fmt.Errorf("query organization: %w", err)
	}

	defer rows.Close()

	organization, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[organization])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.Organization{}, errs.ErrOrganizationNotFound
		}

		return dto.Organization{}, fmt.Errorf("collect organization row: %w", err)
	}

	return convertToApplication(organization), nil
}