package di

import accountRepo "github.com/kirillVladov/account-service/internal/repository/postgres/account"

func (di *DI) AccountRepository() *accountRepo.Repository {
	return accountRepo.New(di.Database())
}
