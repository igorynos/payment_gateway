package consumer

import (
	"context"
	"fmt"

	"payment_gateway/internal/worker"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Client struct {
	client  *kgo.Client
	pending map[messageID]*kgo.Record
}

type messageID struct {
	topic     string
	partition int32
	offset    int64
}

func (c *Client) Fetch(
	ctx context.Context,
) (worker.Message, error) {
	fetches := c.client.PollRecords(ctx, 1)

	if err := fetches.Err(); err != nil {
		return worker.Message{}, fmt.Errorf("poll Kafka: %w", err)
	}

	records := fetches.Records()
	if len(records) == 0 {
		if ctx.Err() != nil {
			return worker.Message{}, ctx.Err()
		}

		return worker.Message{}, fmt.Errorf("Kafka returned no records")
	}

	record := records[0]

	id := messageID{
		topic:     record.Topic,
		partition: record.Partition,
		offset:    record.Offset,
	}

	c.pending[id] = record

	return worker.Message{
		Topic:     record.Topic,
		Key:       record.Key,
		Value:     record.Value,
		Partition: record.Partition,
		Offset:    record.Offset,
	}, nil
}

func (c *Client) Commit(
	ctx context.Context,
	message worker.Message,
) error {
	id := messageID{
		topic:     message.Topic,
		partition: message.Partition,
		offset:    message.Offset,
	}

	record, exists := c.pending[id]
	if !exists {
		return fmt.Errorf(
			"pending Kafka record not found: topic=%s partition=%d offset=%d",
			message.Topic,
			message.Partition,
			message.Offset,
		)
	}

	defer c.client.AllowRebalance()

	if err := c.client.CommitRecords(ctx, record); err != nil {
		return fmt.Errorf("commit Kafka record: %w", err)
	}

	delete(c.pending, id)
	return nil
}

func NewClient(
	brokers []string,
	group string,
	topics []string,
) (*Client, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("Kafka brokers are empty")
	}

	if group == "" {
		return nil, fmt.Errorf("Kafka consumer group is empty")
	}

	if len(topics) == 0 {
		return nil, fmt.Errorf("Kafka topics are empty")
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ClientID("payment-gateway-worker"),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topics...),
		kgo.DisableAutoCommit(),
		kgo.BlockRebalanceOnPoll(),
	)
	if err != nil {
		return nil, fmt.Errorf("create Kafka consumer client: %w", err)
	}

	return &Client{
		client:  client,
		pending: make(map[messageID]*kgo.Record),
	}, nil
}

func (c *Client) Close() {
	if c == nil || c.client == nil {
		return
	}

	c.client.Close()
}
