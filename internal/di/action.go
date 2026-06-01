package di

import (
	"github.com/kirillVladov/account-service/internal/application/action/create_user"
	"github.com/kirillVladov/account-service/internal/application/action/get_by_telegram_id"
	"github.com/kirillVladov/account-service/internal/application/action/get_user"
)

func (di *DI) CreateUserAction() *create_user.CreateUserAction {
	return create_user.New(di.AccountRepository(), di.TokenManager())
}

func (di *DI) GetUserAction() *get_user.GetUserAction {
	return get_user.New(di.AccountRepository())
}

func (di *DI) GetByTelegramIDAction() *get_by_telegram_id.Action {
	return get_by_telegram_id.New(di.AccountRepository())
}
