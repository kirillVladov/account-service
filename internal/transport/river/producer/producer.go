package river_producer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	riverqueue "github.com/riverqueue/river"

	river_transport "github.com/kirillVladov/account-service/internal/transport/river"
	river_client "github.com/kirillVladov/account-service/pkg/river"
	tx_manager "github.com/kirillVladov/account-service/pkg/tx"
)

type Producer struct {
	client *river_client.Client
}

func New(client *river_client.Client) *Producer {
	return &Producer{
		client: client,
	}
}

func (p *Producer) ProduceAccountConfirmationEvent(
	ctx context.Context,
	accountID uuid.UUID,
	organizationID int64,
) error {
	args := river_transport.AccountConfirmationEvent{
		AccountID:      accountID,
		OrganizationID: organizationID,
	}

	if err := p.insert(ctx, args, river_transport.AccountConfirmationQueue); err != nil {
		return fmt.Errorf("insert account confirmation event: %w", err)
	}

	return nil
}

func (p *Producer) insert(ctx context.Context, args riverqueue.JobArgs, queue string) error {
	opts := &riverqueue.InsertOpts{
		Queue: queue,
	}

	var err error

	if tx, ok := tx_manager.TxFromContext(ctx); ok {
		_, err = p.client.GetClient().InsertTx(ctx, tx, args, opts)
	} else {
		_, err = p.client.GetClient().Insert(ctx, args, opts)
	}

	if err != nil {
		return fmt.Errorf("insert job: %w", err)
	}

	return nil
}
