package di

import (
	"time"

	"github.com/kirillVladov/account-service/pkg/token_manager"
)

func (di *DI) TokenManager() *token_manager.Manager {
	return token_manager.New(token_manager.Config{
		Secret:            []byte(di.Config().JwtSecret),
		AccessTTL:         time.Hour * 24,      // todo to config
		RefreshTTL:        time.Hour * 24 * 30, // todo to config
		RefreshTokenBytes: 32,                  // todo to config
	})
}
