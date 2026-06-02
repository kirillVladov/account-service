package get_user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	get_user_mock "github.com/kirillVladov/account-service/internal/application/action/get_user/mocks"
	"github.com/kirillVladov/account-service/internal/application/dto"
)

func TestGetUserAction_Do_success(t *testing.T) {
	repo := new(get_user_mock.AccountRepositoryMock)
	a := New(repo)

	ctx := context.Background()
	id := uuid.New()
	expected := dto.Account{
		ID:    id,
		Email: "test@example.com",
		Name:  "Test",
	}

	repo.On("GetByID", ctx, id).Return(expected, nil).Once()

	actual, err := a.Do(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestGetUserAction_Do_repoError(t *testing.T) {
	repo := new(get_user_mock.AccountRepositoryMock)
	a := New(repo)

	ctx := context.Background()
	id := uuid.New()
	repoErr := errors.New("db error")

	repo.On("GetByID", ctx, id).Return(dto.Account{}, repoErr).Once()

	_, err := a.Do(ctx, id)
	require.Error(t, err)
	require.ErrorIs(t, err, repoErr)
	require.Contains(t, err.Error(), "get account")
}
