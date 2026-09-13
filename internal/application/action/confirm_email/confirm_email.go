package confirm_email

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
	"github.com/kirillVladov/account-service/internal/application/dto/errs"
	"github.com/kirillVladov/account-service/pkg/token_manager"
)

type AccountRepository interface {
	GetByID(ctx context.Context, id uuid.UUID, organizationID int64) (dto.Account, error)
	SetConfirmed(ctx context.Context, id uuid.UUID, organizationID int64) error
}

type TokensRepository interface {
	GetByTokenHash(ctx context.Context, tokenHash string) (dto.AccountToken, error)
}

type Action struct {
	accountRepo AccountRepository
	tokensRepo  TokensRepository
}

func New(accountRepo AccountRepository, tokensRepo TokensRepository) *Action {
	return &Action{
		accountRepo: accountRepo,
		tokensRepo:  tokensRepo,
	}
}

func (a *Action) Confirm(ctx context.Context, confirmationToken string) error {
	tokenHash := token_manager.Hash(confirmationToken)

	token, err := a.tokensRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("get token by hash: %w", err)
	}

	if token.TokenType != dto.TokenTypeEmailConfirm {
		return errs.ErrTokenNotValid
	}

	account, err := a.accountRepo.GetByID(ctx, token.UserID, token.OrganizationID)
	if err != nil {
		return fmt.Errorf("get account: %w", err)
	}

	if account.IsConfirmed {
		return errs.ErrAccountAlreadyConfirmed
	}

	if err := a.accountRepo.SetConfirmed(ctx, account.ID, account.OrganizationID); err != nil {
		return fmt.Errorf("set confirmed: %w", err)
	}

	return nil
}
