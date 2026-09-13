package river_client

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	riverqueue "github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

type Client struct {
	client *riverqueue.Client[pgx.Tx]
}

func NewClient(db *pgxpool.Pool, cfg *riverqueue.Config) (*Client, error) {
	client, err := riverqueue.NewClient(riverpgxv5.New(db), cfg)
	if err != nil {
		return nil, fmt.Errorf("init river client: %w", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) Start(ctx context.Context) error {
	return c.client.Start(ctx)
}

func (c *Client) Stop(ctx context.Context) error {
	return c.client.Stop(ctx)
}

func (c *Client) GetClient() *riverqueue.Client[pgx.Tx] {
	return c.client
}
