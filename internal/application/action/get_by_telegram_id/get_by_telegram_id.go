package get_by_telegram_id

import (
	"context"
	"fmt"

	"github.com/kirillVladov/account-service/internal/application/dto"
)

type Repository interface {
	GetByTelegramID(ctx context.Context, telegramID string) (dto.Account, error)
}

type Action struct {
	repo Repository
}

func New(repo Repository) *Action {
	return &Action{
		repo: repo,
	}
}

func (a *Action) Get(ctx context.Context, tgID string) (dto.Account, error) {
	account, err := a.repo.GetByTelegramID(ctx, tgID)
	if err != nil {
		return dto.Account{}, fmt.Errorf("get user by tg id: %w", err)
	}

	return account, nil
}
