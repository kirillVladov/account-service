package river_consumer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	riverqueue "github.com/riverqueue/river"

	river_transport "github.com/kirillVladov/account-service/internal/transport/river"
)

type Action interface {
	Send(ctx context.Context, accountID uuid.UUID, organizationID int64, confirmationLink string) error
}

type AccountConfirmationWorker struct {
	riverqueue.WorkerDefaults[river_transport.AccountConfirmationEvent]
	action Action
}

func New(action Action) *AccountConfirmationWorker {
	return &AccountConfirmationWorker{
		action: action,
	}
}

func (w *AccountConfirmationWorker) Work(ctx context.Context, job *riverqueue.Job[river_transport.AccountConfirmationEvent]) error {
	if err := w.action.Send(ctx, job.Args.AccountID, job.Args.OrganizationID, job.Args.ConfirmationLink); err != nil {
		return fmt.Errorf("send account confirmation request: %w", err)
	}

	return nil
}
