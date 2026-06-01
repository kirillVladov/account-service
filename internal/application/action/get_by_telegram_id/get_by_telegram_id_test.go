package get_by_telegram_id

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	get_by_telegram_id_mock "github.com/kirillVladov/account-service/internal/application/action/get_by_telegram_id/mocks"
	"github.com/kirillVladov/account-service/internal/application/dto"
)

func TestGetByTelegramIDAction_Get_success(t *testing.T) {
	repo := new(get_by_telegram_id_mock.RepositoryMock)
	a := New(repo)

	ctx := context.Background()
	tgID := "123456789"

	expectedAccount := dto.Account{
		ID:         uuid.New(),
		TelegramID: tgID,
		Email:      "test@example.com",
		Name:       "Test User",
	}

	repo.On("GetByTelegramID", ctx, tgID).Return(expectedAccount, nil).Once()

	account, err := a.Get(ctx, tgID)
	require.NoError(t, err)
	require.Equal(t, expectedAccount, account)
}

func TestGetByTelegramIDAction_Get_repoError(t *testing.T) {
	repo := new(get_by_telegram_id_mock.RepositoryMock)
	a := New(repo)

	ctx := context.Background()
	tgID := "123456789"

	repoErr := errors.New("db error")
	repo.On("GetByTelegramID", ctx, tgID).Return(dto.Account{}, repoErr).Once()

	account, err := a.Get(ctx, tgID)
	require.Error(t, err)
	require.ErrorIs(t, err, repoErr)
	require.Contains(t, err.Error(), "get user by tg id")
	require.Equal(t, dto.Account{}, account)
}