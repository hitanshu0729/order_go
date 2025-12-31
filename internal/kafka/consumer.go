package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hitanshu0729/order_go/internal/domain"
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
	dlqProducer *DLQProducer,
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
			if isPoisonError(err) {
				log.Println("☠️ poison message, sending to DLQ:", err)

				dlqProducer.Publish(ctx, kafka.Message{
					Value: buildDLQPayload(msg, err),
				})

				// ✅ commit offset so partition moves on
				c.reader.CommitMessages(ctx, msg)
				continue
			}

			log.Println("⚠️ transient error, retrying:", err)
			continue // ❌ do not commit offset
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Println("❌ offset commit failed:", err)
		}
	}
}

func isPoisonError(err error) bool {
	return errors.Is(err, domain.ErrInsufficientStock) ||
		errors.Is(err, domain.ErrInvalidPayload) ||
		errors.Is(err, domain.ErrOrderNotFound) ||
		errors.Is(err, domain.ErrProductNotFound)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
func buildDLQPayload(msg kafka.Message, err error) []byte {
	return []byte(
		fmt.Sprintf(
			`{"original_topic":"%s","partition":%d,"offset":%d,"error":"%s","payload":%s}`,
			msg.Topic,
			msg.Partition,
			msg.Offset,
			err.Error(),
			string(msg.Value),
		),
	)
}
