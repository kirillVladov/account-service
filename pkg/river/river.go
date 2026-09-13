package river_client

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	riverqueue "github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

// const EventWorkflow = "event_workflow"

// type Event struct {
// 	ID    uint64         `json:"id" river:"unique"`
// 	Delay *time.Duration `json:"duration"`
// 	Type  string         `json:"type" river:"unique"`
// }

// func (Event) Kind() string {
// 	return "event_workflow"
// }

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

// func (c *Client) Insert(ctx context.Context, event riverqueue.JobArgs) error {
// 	e, ok := event.(Event)
// 	if !ok {
// 		return fmt.Errorf("unknown event type")
// 	}

// 	opts := &riverqueue.InsertOpts{
// 		Queue: EventWorkflow,
// 	}

// 	if e.Delay != nil {
// 		opts.ScheduledAt.Add(*e.Delay)
// 	}

// 	_, err := c.client.Insert(ctx, event, opts)
// 	if err != nil {
// 		return fmt.Errorf("insert job: %w", err)
// 	}

// 	return nil
// }
