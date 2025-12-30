package sqlite

import (
	"context"
	"database/sql"
	"log"

	"github.com/hitanshu0729/order_go/internal/models"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) CreateUser(ctx context.Context, name, email string) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (name, email) VALUES (?, ?)`,
		name,
		email,
	)
	return err
}

func (r *Repo) GetUsers(ctx context.Context) ([]*models.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, email FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			log.Println("Error scanning user:", err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *Repo) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, email FROM users WHERE id = ?`, id)
	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			log.Println("User not found with ID:", id)
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repo) DeleteUser(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	return err
}

func (r *Repo) UpdateUser(ctx context.Context, id int64, name, email string) error {
	user, err := r.db.ExecContext(
		ctx,
		`UPDATE users SET name = ?, email = ? WHERE id = ?`,
		name,
		email,
		id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := user.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("Succesfully Update user with id : %d , num rows affected: %d", id, rowsAffected)
	return nil
}
