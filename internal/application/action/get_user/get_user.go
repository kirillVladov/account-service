package get_user

import (
	"context"

	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
)

type GetUserAction struct {
}

func New() *GetUserAction {
	return &GetUserAction{}
}

func (a *GetUserAction) Do(ctx context.Context, id uuid.UUID) (dto.Account, error) {
	return dto.Account{}, nil
}
