package create_user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
)

type AccountRepository interface {
	Create(ctx context.Context, account dto.Account) error
}

type IssuePair interface {
	IssuePair(userID, role string) (string, string, error)
}

type CreateUserAction struct {
	repo         AccountRepository
	tokenManager IssuePair
}

func New(repo AccountRepository, tokenManager IssuePair) *CreateUserAction {
	return &CreateUserAction{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

func (a *CreateUserAction) Do(ctx context.Context, account dto.Account) error {
	account.ID = uuid.New()

	token, refreshToken, err := a.tokenManager.IssuePair(account.ID.String(), string(dto.UserRoleUser))
	if err != nil {
		return fmt.Errorf("issue token pair: %w", err)
	}

	account.Token = token
	account.RefreshToken = refreshToken

	err = a.repo.Create(ctx, account)
	if err != nil {
		return fmt.Errorf("create account: %w", err)
	}

	return nil
}
