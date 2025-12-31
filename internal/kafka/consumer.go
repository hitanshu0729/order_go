package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          "orders.events",
			GroupID:        groupID,
			MinBytes:       1e3,         // 1KB
			MaxBytes:       10e6,        // 10MB
			CommitInterval: time.Second, // auto-commit
		}),
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("📥 Kafka consumer started")

	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("❌ consumer error:", err)
			return
		}

		log.Printf(
			"📨 event received | topic=%s partition=%d offset=%d value=%s",
			msg.Topic,
			msg.Partition,
			msg.Offset,
			string(msg.Value),
		)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
