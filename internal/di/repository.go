package di

import (
	accountRepo "github.com/kirillVladov/account-service/internal/repository/postgres/account"
	"github.com/kirillVladov/account-service/internal/repository/postgres/account_tokens"
	organizationRepo "github.com/kirillVladov/account-service/internal/repository/postgres/organization"
)

func (di *DI) AccountRepository() *accountRepo.Repository {
	return accountRepo.New(di.Database())
}

func (di *DI) AccountTokenRepository() *account_tokens.Repository {
	return account_tokens.New(di.Database())
}

func (di *DI) OrganizationRepository() *organizationRepo.Repository {
	return organizationRepo.New(di.Database())
}
