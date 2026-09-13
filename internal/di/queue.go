package di

import (
	"github.com/riverqueue/river"

	river_transport "github.com/kirillVladov/account-service/internal/transport/river"
	river_consumer "github.com/kirillVladov/account-service/internal/transport/river/consumer"
	river_client "github.com/kirillVladov/account-service/pkg/river"
)

func (di *DI) AccountConfirmationQueue() *river_client.Client {
	if di.queue != nil {
		return di.queue
	}

	accountConfirmationWorker := river_consumer.New(di.SendMailForConfirmAccount())

	workers := river.NewWorkers()
	river.AddWorker(workers, accountConfirmationWorker)

	cfg := &river.Config{
		Queues: map[string]river.QueueConfig{
			river_transport.AccountConfirmationQueue: {MaxWorkers: 10},
		},
		Workers: workers,
	}

	queue, err := river_client.NewClient(
		di.Database(),
		cfg,
	)
	if err != nil {
		panic("init queue")
	}

	di.queue = queue

	return di.queue
}
