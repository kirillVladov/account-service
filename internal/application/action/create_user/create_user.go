package create_user

import (
	"context"
	"fmt"

	"github.com/kirillVladov/account-service/internal/application/dto"
)

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) (err error)
}

type CreateUserAction struct {
	tx TxManager
}

func New(tx TxManager) *CreateUserAction {
	return &CreateUserAction{
		tx: tx,
	}
}

func (a *CreateUserAction) Do(ctx context.Context, account dto.Account) error {
	err := a.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		// fill data in all tables
		// a.repository.

		return nil
	})
	if err != nil {
		return fmt.Errorf("create account: %w", err)
	}

	return nil
}
