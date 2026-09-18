package di

import (
	"github.com/kirillVladov/account-service/internal/application/action/confirm_email"
	"github.com/kirillVladov/account-service/internal/application/action/create_user"
	"github.com/kirillVladov/account-service/internal/application/action/get_user"
	"github.com/kirillVladov/account-service/internal/application/action/login_user"
	"github.com/kirillVladov/account-service/internal/application/action/queue_email_confirmation"
	refreshtoken_action "github.com/kirillVladov/account-service/internal/application/action/refresh_token"
	"github.com/kirillVladov/account-service/internal/application/action/send_email_confirmation_email"
	notification_gateway "github.com/kirillVladov/account-service/internal/gateway/grpc/notification"
	river_producer "github.com/kirillVladov/account-service/internal/transport/river/producer"
)

func (di *DI) CreateUserAction() *create_user.CreateUserAction {
	return create_user.New(
		di.AccountRepository(),
		di.TokenManager(),
		di.AccountTokenRepository(),
		di.QueueEmailConfirmationAction(),
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
		di.AccountRepository(),
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

func (di *DI) Producer() *river_producer.Producer {
	return river_producer.New(di.AccountConfirmationQueue())
}

func (di *DI) QueueEmailConfirmationAction() *queue_email_confirmation.Action {
	return queue_email_confirmation.New(
		di.AccountRepository(),
		di.AccountTokenRepository(),
		di.Producer(),
		di.TxManager(),
		di.config.AuthToken.EmailConfirmationTokenTTL,
	)
}
