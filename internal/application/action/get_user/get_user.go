package get_user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
)

type AccountRepository interface {
	GetByID(ctx context.Context, id uuid.UUID, organizationID int64) (dto.Account, error)
}

type GetUserAction struct {
	repo AccountRepository
}

func New(repo AccountRepository) *GetUserAction {
	return &GetUserAction{
		repo: repo,
	}
}

func (a *GetUserAction) Do(ctx context.Context, id uuid.UUID, organizationID int64) (dto.Account, error) {
	account, err := a.repo.GetByID(ctx, id, organizationID)
	if err != nil {
		return dto.Account{}, fmt.Errorf("get account: %w", err)
	}

	return account, nil
}
