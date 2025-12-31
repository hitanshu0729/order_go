package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hitanshu0729/order_go/internal/storage/sqlite"
)

type InventoryConsumer struct {
	repo *sqlite.Repo
}

func NewInventoryConsumer(repo *sqlite.Repo) *InventoryConsumer {
	return &InventoryConsumer{repo: repo}
}

type Event struct {
	Type    string       `json:"type"`
	Payload EventPayload `json:"payload"`
}

type EventPayload struct {
	OrderID int64 `json:"order_id"`
}

func (c *InventoryConsumer) HandleMessage(ctx context.Context, value []byte) error {
	var e Event
	if err := json.Unmarshal(value, &e); err != nil {
		return err
	}
	log.Printf("Event received: %+v", e)
	if e.Type != "order.paid" {
		return nil // ignore other events
	}

	log.Println("📦 Inventory update for order:", e.Payload.OrderID)

	return c.processOrder(ctx, e.Payload.OrderID)
}

func (c *InventoryConsumer) processOrder(ctx context.Context, orderID int64) error {
	tx, err := c.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	items, err := c.repo.GetOrderItemsTx(ctx, tx, orderID)
	if err != nil {
		return err
	}

	for _, item := range items {
		err := c.repo.DecreaseProductStockTx(
			ctx,
			tx,
			item.ProductID,
			item.Quantity,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
