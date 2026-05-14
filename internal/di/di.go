package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"

	"github.com/kirillVladov/account-service/internal/config"
)

type DI struct {
	config *config.Config

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
