package kafka

import (
	"context"
	"log"

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
			MinBytes:       1e3,  // 1KB
			MaxBytes:       10e6, // 10MB
			CommitInterval: 0,    // disable auto-commit

		}),
	}
}

func (c *Consumer) Start(
	ctx context.Context,
	inventoryConsumer *InventoryConsumer,
) {
	log.Println("📥 Kafka consumer started")

	for {
		msg, err := c.reader.FetchMessage(ctx)
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

		err = inventoryConsumer.HandleMessage(ctx, msg.Value)
		if err != nil {
			log.Println("inventory update failed:", err)
			// DO NOT crash
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Println("❌ offset commit failed:", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
