package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"

	"github.com/kirillVladov/account-service/internal/config"
	"github.com/kirillVladov/account-service/pkg/notification_client"
	river_client "github.com/kirillVladov/account-service/pkg/river"
)

type DI struct {
	config             *config.Config
	logger             *zap.Logger
	notificatoinClient *notification_client.NotificationClient
	queue              *river_client.Client

	db *pgxpool.Pool
}

func New() *DI {
	return &DI{}
}

func (di *DI) Config() *config.Config {
	if di.config != nil {
		return di.config
	}

	c := config.Config{}

	err := envconfig.Process("", &c)
	if err != nil {
		panic("init config")
	}

	di.config = &c

	return di.config
}

func (di *DI) Logger() *zap.Logger {
	if di.logger != nil {
		return di.logger
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic("init logger")
	}

	di.logger = logger

	return di.logger
}
