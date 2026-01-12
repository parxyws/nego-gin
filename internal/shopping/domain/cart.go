package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameCart = "carts"

// Cart mapped from table <carts>
type Cart struct {
	CartID    string         `gorm:"column:cart_id;primaryKey" json:"cart_id"`
	UserID    string         `gorm:"column:user_id" json:"user_id"`
	SessionID string         `gorm:"column:session_id" json:"session_id"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName Cart's table name
func (*Cart) TableName() string {
	return TableNameCart
}
