package refreshtoken_action

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
	"github.com/kirillVladov/account-service/pkg/token_manager"
)

type TokenManager interface {
	ValidateAccess(raw string) (*token_manager.Claims, error)
	IssuePair(userID, role string) (string, string, error)
}

type AccountTokensRepository interface {
	DeactivateByUser(ctx context.Context, userID uuid.UUID) error
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
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

func (a *RefreshTokenAction) Refresh(ctx context.Context, token string) (string, string, error) {
	claims, err := a.tokenManager.ValidateAccess(token)
	// if err != nil {
	// 	return "", "", fmt.Errorf("validate access: %w", err)
	// }

	token, refreshToken, err := a.tokenManager.IssuePair(claims.UserID, string(dto.UserRoleUser))
	if err != nil {
		return "", "", fmt.Errorf("generate pairs: %w", err)
	}

	err = a.accountTokenRepository.DeactivateByUser(ctx, claims.UserID)
	if err != nil {
		return "", "", fmt.Errorf("deactivate user: %w", err)
	}

	err = a.accountTokenRepository.CreateRefreshToken(ctx, claims.UserID, refreshToken, time.Now())
	if err != nil {
		return "", "", fmt.Errorf("deactivate user: %w", err)
	}

	return token, refreshToken, nil
}
