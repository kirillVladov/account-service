package di

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (di *DI) Database() *pgxpool.Pool {
	if di.db != nil {
		return di.db
	}

	ctx := context.Background()

	cfg, err := pgxpool.ParseConfig(di.Config().PostgresUrl)
	if err != nil {
		panic("parse db cfg")
	}

	// cfg.MaxConns = 20 // todo: turn to cfg
	// cfg.MinConns = 5
	// cfg.MaxConnLifetime = time.Hour
	// cfg.MaxConnIdleTime = 30 * time.Minute
	// cfg.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		panic("apply db cfg")
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		panic("ping db")
	}

	di.db = pool

	return di.db
}
