package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameProduct = "products"

type Product struct {
	ProductID         string         `gorm:"column:product_id;primaryKey" json:"product_id"`
	SellerID          string         `gorm:"column:seller_id;not null" json:"seller_id"`
	CategoryID        int32          `gorm:"column:category_id;not null" json:"category_id"`
	ProductType       string         `gorm:"column:product_type;not null" json:"product_type"`
	Sku               string         `gorm:"column:sku;not null" json:"sku"`
	Name              string         `gorm:"column:name;not null" json:"name"`
	Slug              string         `gorm:"column:slug;not null" json:"slug"`
	Description       string         `gorm:"column:description" json:"description"`
	ShortDescription  string         `gorm:"column:short_description" json:"short_description"`
	BasePrice         float64        `gorm:"column:base_price" json:"base_price"`
	SalePrice         float64        `gorm:"column:sale_price" json:"sale_price"`
	CostPrice         float64        `gorm:"column:cost_price" json:"cost_price"`
	StockQuantity     int32          `gorm:"column:stock_quantity" json:"stock_quantity"`
	LowStockThreshold int32          `gorm:"column:low_stock_threshold;default:5" json:"low_stock_threshold"`
	IsUnlimitedStock  bool           `gorm:"column:is_unlimited_stock;not null" json:"is_unlimited_stock"`
	Status            string         `gorm:"column:status;not null;default:draft" json:"status"`
	IsFeatured        bool           `gorm:"column:is_featured;not null" json:"is_featured"`
	WeightKg          float64        `gorm:"column:weight_kg" json:"weight_kg"`
	DimensionsCm      string         `gorm:"column:dimensions_cm" json:"dimensions_cm"`
	Brand             string         `gorm:"column:brand" json:"brand"`
	Condition         string         `gorm:"column:condition" json:"condition"`
	CreatedAt         time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (*Product) TableName() string {
	return TableNameProduct
}
