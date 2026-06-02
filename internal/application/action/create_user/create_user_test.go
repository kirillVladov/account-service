package create_user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	create_user_mock "github.com/kirillVladov/account-service/internal/application/action/create_user/mocks"
	"github.com/kirillVladov/account-service/internal/application/dto"
)

func TestCreateUserAction_Do_success(t *testing.T) {
	repo := new(create_user_mock.AccountRepositoryMock)
	tokenMgr := new(create_user_mock.IssuePairMock)
	a := New(repo, tokenMgr)

	ctx := context.Background()
	input := dto.Account{
		Email: "test@example.com",
		Name:  "Test",
	}

	// Expect token issuance with generated userID and USER role.
	issuedToken := "token"
	issuedRefresh := "refresh"

	tokenMgr.On("IssuePair", mock.AnythingOfType("string"), string(dto.UserRoleUser)).
		Return(issuedToken, issuedRefresh, nil).
		Once()

	repo.On("Create", ctx, mock.MatchedBy(func(acc dto.Account) bool {
		return acc.ID != (dto.Account{}.ID) &&
			acc.Email == input.Email &&
			acc.Name == input.Name
	})).
		Return(nil).
		Once()

	err := a.Do(ctx, input)
	require.NoError(t, err)
}

func TestCreateUserAction_Do_issuePairError(t *testing.T) {
	repo := new(create_user_mock.AccountRepositoryMock)
	tokenMgr := new(create_user_mock.IssuePairMock)
	a := New(repo, tokenMgr)

	ctx := context.Background()
	input := dto.Account{Email: "test@example.com", Name: "Test"}

	issueErr := errors.New("issuer down")
	tokenMgr.On("IssuePair", mock.AnythingOfType("string"), string(dto.UserRoleUser)).
		Return("", "", issueErr).
		Once()

	err := a.Do(ctx, input)
	require.Error(t, err)
	require.ErrorIs(t, err, issueErr)
	require.Contains(t, err.Error(), "issue token pair")
}

func TestCreateUserAction_Do_repoError(t *testing.T) {
	repo := new(create_user_mock.AccountRepositoryMock)
	tokenMgr := new(create_user_mock.IssuePairMock)
	a := New(repo, tokenMgr)

	ctx := context.Background()
	input := dto.Account{Email: "test@example.com", Name: "Test"}

	repoErr := errors.New("db error")

	tokenMgr.On("IssuePair", mock.AnythingOfType("string"), string(dto.UserRoleUser)).
		Return("token", "refresh", nil).
		Once()

	repo.On("Create", ctx, mock.AnythingOfType("dto.Account")).Return(repoErr).Once()

	err := a.Do(ctx, input)
	require.Error(t, err)
	require.ErrorIs(t, err, repoErr)
	require.Contains(t, err.Error(), "create account")
}
