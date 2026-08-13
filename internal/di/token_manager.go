package di

import (
	"github.com/kirillVladov/account-service/pkg/token_manager"
)

func (di *DI) TokenManager() *token_manager.Manager {
	return token_manager.New(token_manager.Config{
		Secret:            di.config.AuthToken.TokenSecret,
		RefreshTTL:        di.config.AuthToken.RefreshTokenTTL,
		AccessTTL:         di.config.AuthToken.TokenTTL,
		RefreshTokenBytes: int(di.Config().AuthToken.RefreshTokenBytesCount),
	})
}
