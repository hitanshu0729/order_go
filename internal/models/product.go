package models

import _ "gorm.io/gorm"

type Product struct {
	ID    int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"not null" json:"name"`
	Price int64  `gorm:"not null;check:price > 0" json:"price"`
	Stock int64  `gorm:"not null;check:stock >= 0" json:"stock"`
}
