package worker

import "context"

type Handler interface {
	Handle(ctx context.Context, message Message) error
}
