package domain

import (
	"time"
)

const TableNameOrderItem = "order_items"

// OrderItem mapped from table <order_items>
type OrderItem struct {
	OrderItemID int32     `gorm:"column:order_item_id;primaryKey;autoIncrement:true" json:"order_item_id"`
	OrderID     string    `gorm:"column:order_id;not null" json:"order_id"`
	ProductID   string    `gorm:"column:product_id;not null" json:"product_id"`
	VariantID   int32     `gorm:"column:variant_id" json:"variant_id"`
	ProductName string    `gorm:"column:product_name;not null" json:"product_name"`
	Sku         string    `gorm:"column:sku;not null" json:"sku"`
	Quantity    int32     `gorm:"column:quantity;not null" json:"quantity"`
	UnitPrice   float64   `gorm:"column:unit_price;not null" json:"unit_price"`
	Subtotal    float64   `gorm:"column:subtotal;not null" json:"subtotal"`
	TaxAmount   float64   `gorm:"column:tax_amount;not null" json:"tax_amount"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName OrderItem's table name
func (*OrderItem) TableName() string {
	return TableNameOrderItem
}
