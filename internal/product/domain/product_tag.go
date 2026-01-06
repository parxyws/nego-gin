package domain

import "time"

const TableNameProductTag = "product_tags"

type ProductTag struct {
	ProductID string    `gorm:"column:product_id;primaryKey" json:"product_id"`
	TagID     int32     `gorm:"column:tag_id;primaryKey" json:"tag_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (*ProductTag) TableName() string {
	return TableNameProductTag
}
