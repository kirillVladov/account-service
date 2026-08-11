package create_user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
	"github.com/kirillVladov/account-service/pkg/token_manager"
)

type AccountRepository interface {
	Create(ctx context.Context, account dto.Account) error
}

type TokensRepo interface {
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
}

type IssuePair interface {
	IssuePair(userID, role string) (string, string, error)
}

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type CreateUserAction struct {
	repo          AccountRepository
	tokenManager  IssuePair
	tokensRepo    TokensRepo
	tokenDuration time.Duration
	txManager     TxManager
}

func New(
	repo AccountRepository,
	tokenManager IssuePair,
	tokensRepo TokensRepo,
	tokenDuration time.Duration,
	txManager TxManager,
) *CreateUserAction {
	return &CreateUserAction{
		repo:          repo,
		tokenManager:  tokenManager,
		tokensRepo:    tokensRepo,
		tokenDuration: tokenDuration,
		txManager:     txManager,
	}
}

func (a *CreateUserAction) Do(ctx context.Context, account dto.Account) (dto.Account, string, error) {
	account.ID = uuid.New()

	token, refreshToken, err := a.tokenManager.IssuePair(account.ID.String(), string(dto.UserRoleUser))
	if err != nil {
		return dto.Account{}, "", fmt.Errorf("issue token pair: %w", err)
	}

	err = a.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err = a.repo.Create(ctx, account); err != nil {
			return fmt.Errorf("create account: %w", err)
		}

		expires := time.Now().Add(a.tokenDuration)
		tokenHash := token_manager.HashToken(refreshToken)

		if err = a.tokensRepo.CreateRefreshToken(ctx, account.ID, tokenHash, expires); err != nil {
			return fmt.Errorf("create token: %w", err)
		}

		return nil
	})
	if err != nil {
		return dto.Account{}, "", fmt.Errorf("tx error: %w", err)
	}

	return account, token, nil
}
