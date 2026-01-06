package domain

import "time"

const TableNameProductAttribute = "product_attributes"

type ProductAttribute struct {
	AttributeID    int32     `gorm:"column:attribute_id;primaryKey;autoIncrement:true" json:"attribute_id"`
	ProductID      string    `gorm:"column:product_id;not null" json:"product_id"`
	AttributeName  string    `gorm:"column:attribute_name;not null" json:"attribute_name"`
	AttributeValue string    `gorm:"column:attribute_value;not null" json:"attribute_value"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (*ProductAttribute) TableName() string {
	return TableNameProductAttribute
}
