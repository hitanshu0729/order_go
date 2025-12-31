package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "orders.events",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Publish(ctx context.Context, eventType string, payload any) error {
	body, err := json.Marshal(map[string]any{
		"type":      eventType,
		"payload":   payload,
		"timestamp": time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: body,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
