package producer

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Client struct {
	client *kgo.Client
}

func NewClient(brokers []string) (*Client, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ClientID("payment-gateway-producer"),
	)
	if err != nil {
		return nil, fmt.Errorf("create Kafka producer client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Produce(
	ctx context.Context,
	record *kgo.Record,
) error {
	if err := c.client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return fmt.Errorf("produce Kafka record: %w", err)
	}

	return nil
}

func (c *Client) Close() {
	c.client.Close()
}
