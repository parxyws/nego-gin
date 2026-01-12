package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameCategory = "categories"

// Category mapped from table <categories>
type Category struct {
	CategoryID       int32          `gorm:"column:category_id;primaryKey;autoIncrement:true" json:"category_id"`
	ParentCategoryID int32          `gorm:"column:parent_category_id" json:"parent_category_id"`
	CategoryName     string         `gorm:"column:category_name;not null" json:"category_name"`
	Slug             string         `gorm:"column:slug;not null" json:"slug"`
	Description      string         `gorm:"column:description" json:"description"`
	ImageURL         string         `gorm:"column:image_url" json:"image_url"`
	IsActive         bool           `gorm:"column:is_active;not null;default:true" json:"is_active"`
	SortOrder        int32          `gorm:"column:sort_order" json:"sort_order"`
	CreatedAt        time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

// TableName Category's table name
func (*Category) TableName() string {
	return TableNameCategory
}
