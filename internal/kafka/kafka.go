package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Kafka struct {
	client *kgo.Client
}

func New(brokers []string) (*Kafka, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ClientID("payment_gateway"),
	)

	if err != nil {
		return nil, fmt.Errorf("Create kafka error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("ping kafka: %w", err)
	}

	return &Kafka{
		client: client,
	}, nil
}

func (k *Kafka) Close() {
	k.client.Close()
}
