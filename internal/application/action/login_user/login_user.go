package login_user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
	"github.com/kirillVladov/account-service/internal/application/dto/errs"
	"github.com/kirillVladov/account-service/pkg/token_manager"
)

type AccountRepository interface {
	GetByEmail(ctx context.Context, email string, organizationID int64) (dto.Account, error)
}

type TokensRepo interface {
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, organizationID int64, tokenHash string, expiresAt time.Time) error
}

type IssuePair interface {
	IssuePair(userID, role string, organizationID int64) (string, string, error)
}

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type LoginUserAction struct {
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
) *LoginUserAction {
	return &LoginUserAction{
		repo:          repo,
		tokenManager:  tokenManager,
		tokensRepo:    tokensRepo,
		tokenDuration: tokenDuration,
		txManager:     txManager,
	}
}

func (a *LoginUserAction) Do(ctx context.Context, email, password string, organizationID int64) (dto.Account, string, string, error) {
	var (
		outToken        string
		outRefreshToken string
		outAccount      dto.Account
	)

	err := a.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		account, err := a.repo.GetByEmail(ctx, email, organizationID)
		if err != nil {
			return err
		}

		hashedInPassword := token_manager.Hash(password)

		if account.PasswordHash != hashedInPassword {
			return errs.ErrInvalidCredentials
		}

		token, refreshToken, err := a.tokenManager.IssuePair(account.ID.String(), string(dto.UserRoleUser), account.OrganizationID)
		if err != nil {
			return fmt.Errorf("issue token pair: %w", err)
		}

		expires := time.Now().Add(a.tokenDuration)
		tokenHash := token_manager.Hash(refreshToken)

		if err = a.tokensRepo.CreateRefreshToken(ctx, account.ID, account.OrganizationID, tokenHash, expires); err != nil {
			return fmt.Errorf("create token: %w", err)
		}

		outToken = token
		outRefreshToken = refreshToken
		outAccount = account

		return nil
	})
	if err != nil {
		return dto.Account{}, "", "", fmt.Errorf("tx error: %w", err)
	}

	return outAccount, outToken, outRefreshToken, nil
}
