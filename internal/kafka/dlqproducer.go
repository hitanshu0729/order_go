package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type DLQProducer struct {
	writer *kafka.Writer
}

func NewDLQProducer(brokers []string) *DLQProducer {
	return &DLQProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "orders-events.dlq",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *DLQProducer) Publish(ctx context.Context, msg kafka.Message) error {
	return p.writer.WriteMessages(ctx, msg)
}
