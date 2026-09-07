package worker

import (
	"context"
	"fmt"
)

type Source interface {
	Fetch(ctx context.Context) (Message, error)
	Commit(ctx context.Context, message Message) error
}

type Worker struct {
	source     Source
	dispatcher *Dispatcher
}

func New(
	source Source,
	dispatcher *Dispatcher,
) *Worker {
	return &Worker{
		source:     source,
		dispatcher: dispatcher,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		message, err := w.source.Fetch(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			return fmt.Errorf("fetch message: %w", err)
		}

		if err := w.dispatcher.Handle(ctx, message); err != nil {
			return fmt.Errorf(
				"handle topic=%s partition=%d offset=%d: %w",
				message.Topic,
				message.Partition,
				message.Offset,
				err,
			)
		}

		if err := w.source.Commit(ctx, message); err != nil {
			return fmt.Errorf("commit message: %w", err)
		}
	}
}
