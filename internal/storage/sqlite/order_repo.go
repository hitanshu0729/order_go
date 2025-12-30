// UpdateOrderStatus updates the status of an order by id

package sqlite

import (
	"github.com/hitanshu0729/order_go/internal/models"
	"context"
	"database/sql"
	"strings"
	"time"
)

func (r *Repo) CreateOrder(ctx context.Context, userID int64, status string, totalAmount int64) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO orders (user_id, status, total_amount) VALUES (?, ?, ?)`,
		userID,
		status,
		totalAmount,
	)
	return err
}

func (r *Repo) GetOrders(ctx context.Context) ([]*models.Order, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, status, total_amount, created_at FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.TotalAmount, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, nil
}

func (r *Repo) GetOrderByID(ctx context.Context, id int64) (*models.Order, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, status, total_amount, created_at FROM orders WHERE id = ?`, id)
	var o models.Order
	if err := row.Scan(&o.ID, &o.UserID, &o.Status, &o.TotalAmount, &o.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &o, nil
}

func (r *Repo) GetOrdersByStatus(ctx context.Context, status string) ([]*models.Order, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, status, total_amount, created_at FROM orders WHERE status = ?`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.TotalAmount, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, nil
}

func (r *Repo) UpdateOrderStatus(ctx context.Context, id int64, status string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE orders SET status = ? WHERE id = ?`,
		status,
		id,
	)
	return err
}

// UpdateOrderTotal sets the total_amount for a given order.
func (r *Repo) UpdateOrderTotal(ctx context.Context, orderID int64, total int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE orders SET total_amount = ? WHERE id = ?`, total, orderID)
	return err
}

// OrderFilter holds possible filter fields
type OrderFilter struct {
	UserID *int64
	Status *string
	From   *time.Time
	To     *time.Time
}

func (r *Repo) GetOrdersFiltered(ctx context.Context, filter OrderFilter) ([]*models.Order, error) {
	query := "SELECT id, user_id, status, total_amount, created_at FROM orders"
	var args []interface{}
	var conditions []string

	if filter.UserID != nil {
		conditions = append(conditions, "user_id = ?")
		args = append(args, *filter.UserID)
	}
	if filter.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *filter.Status)
	}
	if filter.From != nil {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, filter.From.Format("2006-01-02"))
	}
	if filter.To != nil {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, filter.To.Format("2006-01-02"))
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.TotalAmount, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, nil
}
