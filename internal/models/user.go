package models

import _ "gorm.io/gorm"

// User represents a user in the system.
type User struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"not null" json:"name"`
	Email string `gorm:"not null;unique" json:"email"`
}
