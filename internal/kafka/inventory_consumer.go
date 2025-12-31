package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hitanshu0729/order_go/internal/domain"
	"github.com/hitanshu0729/order_go/internal/storage/sqlite"
	"github.com/mattn/go-sqlite3"
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
		return fmt.Errorf("%w: %v", domain.ErrInvalidPayload, err)
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

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	// 1️⃣ Idempotency guard
	err = c.repo.MarkEventProcessedTx(
		ctx,
		tx,
		"order.paid",
		orderID,
	)
	if err != nil {
		if isUniqueConstraintError(err) {
			// already processed → safe no-op
			log.Printf("Already proceesed for order.paid event type and order id %d", orderID)
			return nil
		}
		return err
	}

	// 2️⃣ Apply side effects
	items, err := c.repo.GetOrderItemsTx(ctx, tx, orderID)
	if err != nil {
		return err
	}

	for _, item := range items {
		log.Printf("Decreasing stock for product %d by %d", item.ProductID, item.Quantity)
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
	if err := tx.Commit(); err != nil {
		return err
	}

	// 3️⃣ Commit both together
	committed = true
	return nil
}
func isUniqueConstraintError(err error) bool {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique ||
			sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey
	}
	return false
}
