package river_consumer

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	riverqueue "github.com/riverqueue/river"
)

type Action interface {
	Send(ctx context.Context, accountID uuid.UUID, organizationID int64) error
}

type AccountConfirmationWorker struct {
	riverqueue.WorkerDefaults[AccountConfirmationEvent]
	action Action
}

func New(action Action) *AccountConfirmationWorker {
	return &AccountConfirmationWorker{
		action: action,
	}
}

func (w *AccountConfirmationWorker) Work(ctx context.Context, job *riverqueue.Job[AccountConfirmationEvent]) error {
	if err := w.action.Send(ctx, job.Args.AccountID, job.Args.OrganizationID); err != nil {
		return fmt.Errorf("send account confirmation request: %w", err)
	}

	return nil
}
