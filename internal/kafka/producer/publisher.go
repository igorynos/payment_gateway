package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

type RecordProducer interface {
	Produce(
		ctx context.Context,
		record *kgo.Record,
	) error
}

type Publisher struct {
	producer RecordProducer
	topic    string
}

func (p *Publisher) Publish(
	ctx context.Context,
	key string,
	message any,
) error {
	value, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	record := &kgo.Record{
		Topic: p.topic,
		Key:   []byte(key),
		Value: value,
	}

	if err := p.producer.Produce(ctx, record); err != nil {
		return fmt.Errorf("publish message to topic %q: %w", p.topic, err)
	}

	return nil
}

func NewPublisher(
	producer RecordProducer,
	topic string,
) *Publisher {
	return &Publisher{
		producer: producer,
		topic:    topic,
	}
}
