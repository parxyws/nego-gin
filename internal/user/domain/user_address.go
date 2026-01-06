package domain

import (
	"time"

	"gorm.io/gorm"
)

const TableNameUserAddress = "user_addresses"

type UserAddress struct {
	AddressID     int32          `gorm:"column:address_id;primaryKey;autoIncrement:true" json:"address_id"`
	UserID        string         `gorm:"column:user_id;not null" json:"user_id"`
	AddressType   string         `gorm:"column:address_type;not null" json:"address_type"`
	IsDefault     bool           `gorm:"column:is_default;not null" json:"is_default"`
	RecipientName string         `gorm:"column:recipient_name;not null" json:"recipient_name"`
	AddressLine1  string         `gorm:"column:address_line1;not null" json:"address_line1"`
	AddressLine2  string         `gorm:"column:address_line2" json:"address_line2"`
	City          string         `gorm:"column:city;not null" json:"city"`
	StateProvince string         `gorm:"column:state_province" json:"state_province"`
	PostalCode    string         `gorm:"column:postal_code;not null" json:"postal_code"`
	CountryCode   string         `gorm:"column:country_code;not null" json:"country_code"`
	Phone         string         `gorm:"column:phone" json:"phone"`
	CreatedAt     time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (*UserAddress) TableName() string {
	return TableNameUserAddress
}
