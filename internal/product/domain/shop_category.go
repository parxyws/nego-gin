package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameShopCategory = "shop_categories"

// ShopCategory mapped from table <shop_categories>
type ShopCategory struct {
	ShopCategoryID       int32          `gorm:"column:shop_category_id;primaryKey;autoIncrement:true" json:"shop_category_id"`
	SellerID             string         `gorm:"column:seller_id;not null" json:"seller_id"`
	ParentShopCategoryID int32          `gorm:"column:parent_shop_category_id" json:"parent_shop_category_id"`
	CategoryName         string         `gorm:"column:category_name;not null" json:"category_name"`
	Slug                 string         `gorm:"column:slug;not null" json:"slug"`
	Description          string         `gorm:"column:description" json:"description"`
	IsActive             bool           `gorm:"column:is_active;not null;default:true" json:"is_active"`
	SortOrder            int32          `gorm:"column:sort_order" json:"sort_order"`
	CreatedAt            time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	Products             []Product      `gorm:"foreignKey:shop_category_id;references:id" json:"products"`
}

// TableName ShopCategory's table name
func (*ShopCategory) TableName() string {
	return TableNameShopCategory
}
