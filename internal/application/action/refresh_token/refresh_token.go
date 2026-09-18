package refreshtoken_action

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
	"github.com/kirillVladov/account-service/internal/application/dto/errs"
	"github.com/kirillVladov/account-service/pkg/token_manager"
)

type AccountRepository interface {
	GetByID(ctx context.Context, id uuid.UUID, organizationID int64) (dto.Account, error)
}

type TokenManager interface {
	ValidateAccess(raw string) (*token_manager.Claims, error)
	IssuePair(userID, role string, organizationID int64) (string, string, error)
	IssueAccess(userID string, organizationID int64) (string, error)
}

type AccountTokensRepository interface {
	DeactivateByUser(ctx context.Context, userID uuid.UUID, organizationID int64) error
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, organizationID int64, tokenHash string, expiresAt time.Time) error
	GetTokenByUserID(ctx context.Context, userID uuid.UUID, organizationID int64) (dto.AccountToken, error)
}

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type RefreshTokenAction struct {
	tokenManager           TokenManager
	accountTokenRepository AccountTokensRepository
	accountRepository      AccountRepository
	txManager              TxManager
}

func New(tokenManager TokenManager, accountTokenRepository AccountTokensRepository, accountRepository AccountRepository, txManager TxManager) *RefreshTokenAction {
	return &RefreshTokenAction{
		tokenManager:           tokenManager,
		accountTokenRepository: accountTokenRepository,
		accountRepository:      accountRepository,
		txManager:              txManager,
	}
}

func (a *RefreshTokenAction) getAccount(ctx context.Context, userID uuid.UUID, organizationID int64) (dto.Account, error) {
	account, err := a.accountRepository.GetByID(ctx, userID, organizationID)
	if err != nil {
		return dto.Account{}, fmt.Errorf("get account: %w", err)
	}

	if account.IsBlocked {
		return dto.Account{}, errs.ErrAccountBlocked
	}

	return account, nil
}

func (a *RefreshTokenAction) Refresh(ctx context.Context, oldToken, oldRefreshToken string) (string, string, error) {
	claims, err := a.tokenManager.ValidateAccess(oldToken)
	if err == nil {
		return oldToken, oldRefreshToken, nil
	}

	if !errors.Is(err, token_manager.ErrTokenExpired) {
		return "", "", fmt.Errorf("validate access: %w", err)
	}

	userId, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", "", fmt.Errorf("parse userID: %w", err)
	}

	account, err := a.getAccount(ctx, userId, claims.OrganizationID)
	if err != nil {
		return "", "", err
	}

	accountCreds, err := a.accountTokenRepository.GetTokenByUserID(ctx, account.ID, claims.OrganizationID)
	if err != nil {
		return "", "", fmt.Errorf("get account creds: %w", err)
	}

	hashedOldhRefreshToken := token_manager.Hash(oldRefreshToken)

	if accountCreds.Revoked || hashedOldhRefreshToken != accountCreds.TokenHash {
		return "", "", errs.ErrForbidden
	}

	if accountCreds.ExpiresAt.After(time.Now()) {
		token, err := a.tokenManager.IssueAccess(claims.UserID, claims.OrganizationID)
		if err != nil {
			return "", "", fmt.Errorf("generate access token: %w", err)
		}

		return token, oldRefreshToken, nil
	}

	token, refreshToken, err := a.tokenManager.IssuePair(claims.UserID, string(dto.UserRoleUser), claims.OrganizationID)
	if err != nil {
		return "", "", fmt.Errorf("generate pairs: %w", err)
	}

	err = a.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err = a.accountTokenRepository.DeactivateByUser(ctx, account.ID, claims.OrganizationID); err != nil {
			return fmt.Errorf("deactivate user: %w", err)
		}

		if err = a.accountTokenRepository.CreateRefreshToken(ctx, account.ID, claims.OrganizationID, refreshToken, time.Now()); err != nil {
			return fmt.Errorf("create refresh token: %w", err)
		}

		return nil
	})
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}
