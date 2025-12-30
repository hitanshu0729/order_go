package sqlite

import (
	"github.com/hitanshu0729/order_go/internal/models"
	"context"
	"database/sql"
	"errors"
)

// GetOrderItems returns all items for a given order.
func (r *Repo) GetOrderItems(ctx context.Context, orderID int64) ([]*models.OrderItem, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = ?`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// AddOrderItem adds a new item to an order (only if order is pending/created).
func (r *Repo) AddOrderItem(ctx context.Context, orderID, productID, quantity, price int64) error {
	// Check order status
	order, err := r.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil || order.Status != "pending" {
		return errors.New("can only add items to orders with status 'pending'")
	}
	_, err = r.db.ExecContext(ctx,
		`INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`,
		orderID, productID, quantity, price,
	)
	if err != nil {
		return err
	}
	return r.recalculateOrderTotal(ctx, orderID)
}

// UpdateOrderItemQuantity updates the quantity of an item (only if order is pending/created).
func (r *Repo) UpdateOrderItemQuantity(ctx context.Context, orderID, itemID, quantity int64) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE order_items SET quantity = ? WHERE product_id = ? AND order_id = ?`,
		quantity, itemID, orderID,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return r.recalculateOrderTotal(ctx, orderID)
}

// RemoveOrderItem deletes an item from an order (only if order is pending/created).
func (r *Repo) RemoveOrderItem(ctx context.Context, orderID, itemID int64) error {
	order, err := r.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil || order.Status != "pending" {
		return errors.New("can only remove items from orders with status 'pending'")
	}
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM order_items WHERE product_id = ? AND order_id = ?`,
		itemID, orderID,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return r.recalculateOrderTotal(ctx, orderID)
}

// recalculateOrderTotal updates the total_amount in orders table after item changes.
func (r *Repo) recalculateOrderTotal(ctx context.Context, orderID int64) error {
	row := r.db.QueryRowContext(ctx, `SELECT SUM(quantity * price) FROM order_items WHERE order_id = ?`, orderID)
	var total sql.NullInt64
	if err := row.Scan(&total); err != nil {
		return err
	}
	return r.UpdateOrderTotal(ctx, orderID, total.Int64)
}
