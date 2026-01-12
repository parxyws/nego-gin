package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameProductVariant = "product_variants"

type ProductVariant struct {
	VariantID       int32          `gorm:"column:variant_id;primaryKey;autoIncrement:true" json:"variant_id"`
	ProductID       string         `gorm:"column:product_id;not null" json:"product_id"`
	Sku             string         `gorm:"column:sku;not null" json:"sku"`
	VariantName     string         `gorm:"column:variant_name;not null" json:"variant_name"`
	PriceAdjustment float64        `gorm:"column:price_adjustment" json:"price_adjustment"`
	Attributes      string         `gorm:"column:attributes;not null" json:"attributes"`
	StockQuantity   int32          `gorm:"column:stock_quantity" json:"stock_quantity"`
	IsActive        bool           `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt       time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (*ProductVariant) TableName() string {
	return TableNameProductVariant
}
