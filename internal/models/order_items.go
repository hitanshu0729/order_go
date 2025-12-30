package models

// OrderItem represents an item in an order.
type OrderItem struct {
    ID        int64 `gorm:"primaryKey;autoIncrement" json:"id"`
    OrderID   int64 `gorm:"not null;index" json:"order_id"`
    ProductID int64 `gorm:"not null;index" json:"product_id"`
    Quantity  int64 `gorm:"not null;check:quantity > 0" json:"quantity"`
    Price     int64 `gorm:"not null;check:price > 0" json:"price"`
}