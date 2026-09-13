package di

import (
	"github.com/kirillVladov/account-service/internal/application/action/confirm_email"
	"github.com/kirillVladov/account-service/internal/application/action/create_user"
	"github.com/kirillVladov/account-service/internal/application/action/get_user"
	"github.com/kirillVladov/account-service/internal/application/action/login_user"
	refreshtoken_action "github.com/kirillVladov/account-service/internal/application/action/refresh_token"
	"github.com/kirillVladov/account-service/internal/application/action/send_email_confirmation_email"
	notification_gateway "github.com/kirillVladov/account-service/internal/gateway/grpc/notification"
)

func (di *DI) CreateUserAction() *create_user.CreateUserAction {
	return create_user.New(
		di.AccountRepository(),
		di.TokenManager(),
		di.AccountTokenRepository(),
		di.config.AuthToken.RefreshTokenTTL,
		di.TxManager(),
	)
}

func (di *DI) GetUserAction() *get_user.GetUserAction {
	return get_user.New(di.AccountRepository())
}

func (di *DI) RefreshTokenAction() *refreshtoken_action.RefreshTokenAction {
	return refreshtoken_action.New(
		di.TokenManager(),
		di.AccountTokenRepository(),
	)
}

func (di *DI) LoginUserAction() *login_user.LoginUserAction {
	return login_user.New(
		di.AccountRepository(),
		di.TokenManager(),
		di.AccountTokenRepository(),
		di.config.AuthToken.RefreshTokenTTL,
		di.TxManager(),
	)
}

func (di *DI) NotificationGateway() *notification_gateway.Gateway {
	return notification_gateway.New(di.NotificationServiceClient())
}

func (di *DI) SendMailForConfirmAccount() *send_email_confirmation_email.Action {
	return send_email_confirmation_email.New(
		di.AccountRepository(),
		di.NotificationGateway(),
	)
}

func (di *DI) ConfirmEmailAction() *confirm_email.Action {
	return confirm_email.New(
		di.AccountRepository(),
		di.AccountTokenRepository(),
	)
}
