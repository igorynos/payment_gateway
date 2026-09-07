package worker

import (
	"context"
	"fmt"
)

type Dispatcher struct {
	handlers map[string]Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]Handler),
	}
}

func (d *Dispatcher) Register(
	topic string,
	handler Handler,
) error {
	if topic == "" {
		return fmt.Errorf("topic is empty")
	}

	if handler == nil {
		return fmt.Errorf("handler for topic %q is nil", topic)
	}

	if _, exists := d.handlers[topic]; exists {
		return fmt.Errorf(
			"handler for topic %q is already registered",
			topic,
		)
	}

	d.handlers[topic] = handler
	return nil
}

func (d *Dispatcher) Handle(
	ctx context.Context,
	message Message,
) error {
	handler, exists := d.handlers[message.Topic]
	if !exists {
		return fmt.Errorf(
			"handler for topic %q is not registered",
			message.Topic,
		)
	}

	return handler.Handle(ctx, message)
}
