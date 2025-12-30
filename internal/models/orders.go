package models

import "time"

// Order represents an order in the system.
type Order struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int64     `gorm:"not null;index" json:"user_id"`
	Status      string    `gorm:"not null;check:status IN ('pending','paid','cancelled','completed')" json:"status"`
	TotalAmount int64     `gorm:"not null;check:total_amount > 0" json:"total_amount"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
}
