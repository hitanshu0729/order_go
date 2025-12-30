package sqlite

import (
	"context"
	"database/sql"
	"log"

	"github.com/hitanshu0729/order_go/internal/models"
)

// Create inserts a new product
func (r *Repo) CreateProduct(
	ctx context.Context,
	name string,
	price, stock int64,
) error {
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO products (name, price, stock) VALUES (?, ?, ?)`,
		name,
		price,
		stock,
	)
	if err != nil {
		log.Printf("failed to create product: name=%s, price=%d, stock=%d, error=%v", name, price, stock, err)
		return err
	}
	if rows, err := res.RowsAffected(); err == nil {
		log.Printf("successfully created product: name=%s, price=%d, stock=%d, rows=%d", name, price, stock, rows)
	}
	return nil
}

// GetProducts returns all products
func (r *Repo) GetProducts(ctx context.Context) ([]*models.Product, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, name, price, stock FROM products`,
	)
	if err != nil {
		log.Printf("failed to get products: %v", err)
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			log.Printf("error scanning product: %v", err)
			return nil, err
		}
		products = append(products, &p)
	}
	log.Printf("retrieved %d products", len(products))
	return products, nil
}

// GetProductByID returns a product by id
func (r *Repo) GetProductByID(
	ctx context.Context,
	id int64,
) (*models.Product, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, price, stock FROM products WHERE id = ?`,
		id,
	)

	var p models.Product
	if err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("product not found with id: %d", id)
			return nil, nil
		}
		log.Printf("failed to get product by id=%d: %v", id, err)
		return nil, err
	}
	log.Printf("retrieved product: id=%d, name=%s", p.ID, p.Name)
	return &p, nil
}

// Delete deletes a product by id
func (r *Repo) Delete(ctx context.Context, id int64) error {

	res, err := r.db.ExecContext(
		ctx,
		`DELETE FROM products WHERE id = ?`,
		id,
	)
	if err != nil {
		log.Printf("failed to delete product id=%d: %v", id, err)
		return err
	}
	if rows, err := res.RowsAffected(); err == nil {
		log.Printf("successfully deleted product id=%d, rows=%d", id, rows)
	} else {
		log.Printf("deleted product id=%d but failed to get rows affected: %v", id, err)
	}
	return nil
}
