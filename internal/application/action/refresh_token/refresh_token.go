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

type TokenManager interface {
	ValidateAccess(raw string) (*token_manager.Claims, error)
	IssuePair(userID, role string) (string, string, error)
	IssueAccess(userID string) (string, error)
}

type AccountTokensRepository interface {
	DeactivateByUser(ctx context.Context, userID uuid.UUID) error
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	GetTokenByUserID(ctx context.Context, userID uuid.UUID) (dto.AccountToken, error)
}

type RefreshTokenAction struct {
	tokenManager           TokenManager
	accountTokenRepository AccountTokensRepository
}

func New(tokenManager TokenManager, accountTokenRepository AccountTokensRepository) *RefreshTokenAction {
	return &RefreshTokenAction{
		tokenManager:           tokenManager,
		accountTokenRepository: accountTokenRepository,
	}
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
		return "", "", fmt.Errorf("parse user id: %w", err)
	}

	accountCreds, err := a.accountTokenRepository.GetTokenByUserID(ctx, userId)
	if err != nil {
		return "", "", fmt.Errorf("get account creds: %w", err)
	}

	hashedOldhRefreshToken := token_manager.HashToken(oldRefreshToken)

	if accountCreds.Revoked || hashedOldhRefreshToken != accountCreds.TokenHash {
		return "", "", errs.ErrForbidden
	}

	if accountCreds.ExpiresAt.After(time.Now()) {
		token, err := a.tokenManager.IssueAccess(claims.UserID)
		if err != nil {
			return "", "", fmt.Errorf("generate access token: %w", err)
		}

		return token, oldRefreshToken, nil
	}

	token, refreshToken, err := a.tokenManager.IssuePair(claims.UserID, string(dto.UserRoleUser))
	if err != nil {
		return "", "", fmt.Errorf("generate pairs: %w", err)
	}

	err = a.accountTokenRepository.DeactivateByUser(ctx, userId)
	if err != nil {
		return "", "", fmt.Errorf("deactivate user: %w", err)
	}

	err = a.accountTokenRepository.CreateRefreshToken(ctx, userId, refreshToken, time.Now())
	if err != nil {
		return "", "", fmt.Errorf("deactivate user: %w", err)
	}

	return token, refreshToken, nil
}
