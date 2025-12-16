package models

import "time"

// Inventory represents an inventory item
type Inventory struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity" gorm:"default:0"`
	Price       float64   `json:"price" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for Inventory
func (Inventory) TableName() string {
	return "inventories"
}
